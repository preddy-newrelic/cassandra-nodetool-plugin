package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
)

func populateMetrics(ms *metric.MetricSet) error {
	cmdExecutable := strings.TrimSpace(args.Cmd)
	checkNodetoolExists(cmdExecutable)
	localHostID := ""

	infoStr, err := runCommand(cmdExecutable, "info")
	if err != nil {
		log.Error("%s command failed with %s\n", cmdExecutable, err)
		ms.SetMetric("status", -1, metric.GAUGE)
		ms.SetMetric("state", -1, metric.GAUGE)
		return nil
	}
	infoLines := strings.Split(infoStr, "\n")
	for _, infoLine := range infoLines {
		infoLinePieces := strings.Split(infoLine, ":")
		if len(infoLinePieces) > 1 {
			infoLineKey := infoLinePieces[0]
			infoLineKey = strings.TrimSpace(infoLineKey)
			infoLineValue := infoLinePieces[1]
			infoLineValue = strings.TrimSpace(infoLineValue)
			//fmt.Println("X" + infoLineKey + "X")
			if infoLineKey == "ID" {
				localHostID = infoLineValue
			} else {
				if strings.Contains(infoLineKey, "Gossip") {
					ms.SetMetric(infoLineKey, infoLineValue, metric.ATTRIBUTE)
				} else if strings.Contains(infoLineKey, "Exceptions") {
					exceptions, err := strconv.ParseFloat(infoLineValue, 64)
					if err == nil {
						ms.SetMetric(infoLineKey, exceptions, metric.GAUGE)
					}
				}
			}
		}
	}

	log.Debug("local host id is [%s]", localHostID)
	if localHostID == "" {
		ms.SetMetric("status", -1, metric.GAUGE)
		ms.SetMetric("state", -1, metric.GAUGE)
		return nil
	}

	outStr, err := runCommand(cmdExecutable, "status")
	if err != nil {
		log.Error("%s command failed with %s\n", cmdExecutable, err)
		ms.SetMetric("status", -1, metric.GAUGE)
		ms.SetMetric("state", -1, metric.GAUGE)
		return nil
	}
	temp := strings.Split(outStr, "\n")
	for _, line := range temp {
		splitedLine := strings.Fields(line)
		if len(splitedLine) < 8 {
			continue
		}
		if splitedLine[6] != localHostID {
			continue
		}
		status := splitedLine[0]

		pattern := regexp.MustCompile(`([UD])([NLJM])`)
		matched := pattern.FindStringSubmatch(status)

		if len(matched) > 2 {
			upDownStatus := matched[1]
			switch upDownStatus {
			case "U":
				ms.SetMetric("status", 2, metric.GAUGE)
			case "D":
				ms.SetMetric("status", 1, metric.GAUGE)
			default:
				log.Error("Unknown status ", upDownStatus)
			}

			state := matched[2]
			switch state {
			case "N":
				ms.SetMetric("state", 1, metric.GAUGE)
			case "L":
				ms.SetMetric("state", 2, metric.GAUGE)
			case "J":
				ms.SetMetric("state", 3, metric.GAUGE)
			case "M":
				ms.SetMetric("state", 4, metric.GAUGE)
			default:
				log.Error("Unknown state ", state)
			}

		}

		ms.SetMetric("address", splitedLine[1], metric.ATTRIBUTE)

		loadFloat, err := strconv.ParseFloat(splitedLine[2], 64)
		if err == nil {
			ms.SetMetric("load", loadFloat, metric.GAUGE)
		}

		tokensFloat, err := strconv.ParseFloat(splitedLine[4], 64)
		if err == nil {
			ms.SetMetric("tokens", tokensFloat, metric.GAUGE)
		}

		ownsStr := splitedLine[5]
		ownsStr = ownsStr[0 : len(ownsStr)-1]
		ownsFloat, err := strconv.ParseFloat(ownsStr, 64)
		if err == nil {
			ms.SetMetric("owns", ownsFloat, metric.GAUGE)
		}

		ms.SetMetric("hostid", splitedLine[6], metric.ATTRIBUTE)
		ms.SetMetric("rack", splitedLine[7], metric.ATTRIBUTE)
		return nil
	}
	return nil
}

func checkNodetoolExists(cmdExecutable string) {
	path, err := exec.LookPath(cmdExecutable)
	if err != nil {
		if args.Verbose {
			fmt.Printf("%s executable not found in PATH\n", cmdExecutable)
		}
	} else {
		if args.Verbose {
			fmt.Printf("%s executable is in '%s'\n", cmdExecutable, path)
		}
	}
}

func runCommand(cmdExecutable string, nodetoolCommand string) (string, error) {
	cmd := exec.Command(cmdExecutable, nodetoolCommand)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if errStr != "" {
		if args.Verbose {
			fmt.Printf("Errors running command:\n%s\n", errStr)
		}
	}
	return outStr, nil
}

func asValue(value string) interface{} {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}
	return value
}
