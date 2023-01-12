import {DataSourceInstanceSettings, CoreApp} from '@grafana/data'
import {DataSourceWithBackend} from '@grafana/runtime'

import {SQLQuery, FlightSQLDataSourceOptions, DEFAULT_QUERY} from './types'

// import { getBackendSrv } from "@grafana/runtime"

export class FlightSQLDataSource extends DataSourceWithBackend<SQLQuery, FlightSQLDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<FlightSQLDataSourceOptions>) {
    super(instanceSettings)
  }

  getDefaultQuery(_: CoreApp): Partial<SQLQuery> {
    return DEFAULT_QUERY
  }

  // can use this to modify the raw query before sending it to the backend
  // applyTemplateVariables(query: BasicQuery, scopedVars: ScopedVars): Record<string, any> {
  //   return {
  //     ...query,
  //     rawQuery: getTemplateSrv().replace(query.rawQuery, scopedVars),
  //   };
  // }

  getTables(): Promise<any> {
    return this.getResource('/get-tables')
  }

  getColumns(table: string): Promise<any> {
    return this.getResource(`/get-columns?table=${table}`)
  }
}
