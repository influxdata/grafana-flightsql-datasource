import {DataSourcePlugin} from '@grafana/data'
import {FlightSQLDatasource} from './Datasource'
import {ConfigEditor} from './components/ConfigEditor'
import {QueryEditor} from './components/QueryEditor'
import {SQLQuery, FlightSQLDataSourceOptions} from './types'

export const plugin = new DataSourcePlugin<FlightSQLDatasource, SQLQuery, FlightSQLDataSourceOptions>(
  FlightSQLDatasource
)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
