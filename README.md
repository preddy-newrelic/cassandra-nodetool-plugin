# New Relic Infrastructure Integration for cassandra-status

Reports status and metrics for cassandra-status service

## Requirements

### cassandra nodetool
This plugin invokes the cassandra nodetool utility to gather metrics.

### New Relic Infrastructure Agent
This is the description about how to run the Cassandra Status Integration with New Relic Infrastructure agent, so it is required to have the agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).


## Installation

* Download an archive file for the Cassandra-Status Integration
* Place the executables under `bin` directory and the definition file `cassandra-status-definition.yml` in `/var/db/newrelic-infra/newrelic-integrations`
* Set execution permissions for the binary file `nr-cassandra-status`
* Place the integration configuration file `cassandra-status-config.yml.sample` in `/etc/newrelic-infra/integrations.d`

## Configuration

In order to use the Cassandra-Status Integration it is required to configure `cassandra-status-config.yml.sample` file. Firstly, rename the file to `cassandra-status-config.yml`. Then, depending on your needs, specify all instances that you want to monitor. Once this is done, restart the Infrastructure agent.

You can view your data in Insights by creating your own custom NRQL queries. To
do so use **CassandraStatusSample** event types.

The following attributes are reported always
* Status: 
     0 (not reported) if there is an error executing `nodetool info` or `nodetool status` commands
     1 if the commands execute succesfully and return a status of "Down"
     2 if the commands execute successfully and return a status of "Up"

* State :
     0 (not reported) if there is an error executing `nodetool info` or `nodetool status` commands
     1 if the commands execute successfully and return a state of `N`
     2 if the commands execute successfully and return a state of `L`
     3 if the commands execute successfully and return a state of `J`
     4 if the commands execute successfully and return a state of `M`

If the command execute succesfully, then the following attributes are also parsed, They will not be reported however if the commands fails for any reason,

* address
* load
* tokens
* hostid
* owns
* rack




## Compatibility

* Supported OS: linux
* cassandra-status versions:
* Edition:

