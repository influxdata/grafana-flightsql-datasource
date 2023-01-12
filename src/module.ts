import {DataSourcePlugin} from '@grafana/data'
import {FlightSQLDataSource} from './datasource'
import {ConfigEditor} from './components/ConfigEditor'
import {QueryEditor} from './components/QueryEditor'
import {SQLQuery, FlightSQLDataSourceOptions} from './types'

export const plugin = new DataSourcePlugin<FlightSQLDataSource, SQLQuery, FlightSQLDataSourceOptions>(
  FlightSQLDataSource
)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
