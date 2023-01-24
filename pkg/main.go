package main

import (
	"os"

	"github.com/grafana/flight-datasource/pkg/flightsql"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	if err := datasource.Manage("influxdata-flightsql-datasource", flightsql.NewDatasource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
