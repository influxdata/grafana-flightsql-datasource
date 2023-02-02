import {DataSourceInstanceSettings, CoreApp, ScopedVars, VariableWithMultiSupport} from '@grafana/data'
import {DataSourceWithBackend, getTemplateSrv} from '@grafana/runtime'
import {SQLQuery, FlightSQLDataSourceOptions, DEFAULT_QUERY} from './types'

export class FlightSQLDataSource extends DataSourceWithBackend<SQLQuery, FlightSQLDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<FlightSQLDataSourceOptions>) {
    super(instanceSettings)
  }

  getDefaultQuery(_: CoreApp): Partial<SQLQuery> {
    return DEFAULT_QUERY
  }

  quoteLiteral(value: string) {
    return "'" + value.replace(/'/g, "''") + "'"
  }

  interpolateVariable = (value: string | string[] | number, variable: VariableWithMultiSupport) => {
    if (typeof value === 'string') {
      if (variable?.multi || variable?.includeAll) {
        return this.quoteLiteral(value)
      } else {
        return String(value).replace(/'/g, "''")
      }
    }

    if (typeof value === 'number') {
      return value
    }

    if (Array.isArray(value)) {
      const quotedValues = value.map((v) => this.quoteLiteral(v))
      return quotedValues.join(',')
    }

    return value
  }

  applyTemplateVariables(query: SQLQuery, scopedVars: ScopedVars): Record<string, any> {
    const interpolatedQuery: SQLQuery = {
      ...query,
      queryText: getTemplateSrv().replace(query.queryText, scopedVars, this.interpolateVariable),
    }
    return interpolatedQuery
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
