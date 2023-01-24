import {DataSourceInstanceSettings, CoreApp} from '@grafana/data'
import {DataSourceWithBackend} from '@grafana/runtime'
import {SQLQuery, FlightSQLDataSourceOptions, DEFAULT_QUERY} from './types'

export class FlightSQLDataSource extends DataSourceWithBackend<SQLQuery, FlightSQLDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<FlightSQLDataSourceOptions>) {
    super(instanceSettings)
  }

  getDefaultQuery(_: CoreApp): Partial<SQLQuery> {
    return DEFAULT_QUERY
  }

  getSQLInfo(): Promise<any> {
    return this.getResource('/flightsql/sql-info')
  }

  getTables(): Promise<any> {
    return this.getResource('/flightsql/tables')
  }

  getColumns(table: string): Promise<any> {
    return this.getResource(`/flightsql/columns?table=${table}`)
  }

  getMacros(): Promise<any> {
    return this.getResource('/plugin/macros')
  }
}
