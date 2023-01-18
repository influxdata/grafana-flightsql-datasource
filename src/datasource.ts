import {DataSourceInstanceSettings, CoreApp} from '@grafana/data'
import {DataSourceWithBackend} from '@grafana/runtime'

import {SQLQuery, FlightSQLDataSourceOptions, DEFAULT_QUERY} from './types'

// import { getBackendSrv } from "@grafana/runtime"

export class FlightSQLDatasource extends DataSourceWithBackend<SQLQuery, FlightSQLDataSourceOptions> {
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
  getSQLInfo(): Promise<any> {
    return this.getResource('/flightsql/sql-info')
  }

  getTables(): Promise<any> {
    return this.getResource('/flightsql/tables')
  }

  getColumns(table: string): Promise<any> {
    return this.getResource(`/flightsql/columns?table=${table}`)
  }
}
