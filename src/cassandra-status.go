package main

import (
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Cmd string `default:"nodetool" help:"nodetool executable"`
}

const (
	integrationName    = "com.newrelic.cassandra-status"
	integrationVersion = "1.0.1"
)

var args argumentList

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)

	if args.All || args.Metrics {
		ms := integration.NewMetricSet("CassandraNodeSample")
		fatalIfErr(populateMetrics(ms))
	}
	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
