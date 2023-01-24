import {DataQuery, DataSourceJsonData} from '@grafana/data'
import {formatSQL} from './components/sqlFormatter'

export interface SQLQuery extends DataQuery {
  queryText?: string
}

export const DEFAULT_QUERY: Partial<SQLQuery> = {}

/**
 * These are options configured for each DataSource instance
 */
export interface FlightSQLDataSourceOptions extends DataSourceJsonData {
  host?: string
  database?: string
  token?: string
  secure?: boolean
  username?: string
  password?: string
  selectedAuthType?: string
}

export type TablesResponse = {
  tables: string[]
}

export type ColumnsResponse = {
  columns: string[]
}

export const authTypeOptions = [
  {key: 0, label: 'none', title: 'none'},
  {key: 1, label: 'username/password', title: 'username/password'},
  {key: 2, label: 'token', title: 'token'},
]

export const sqlLanguageDefinition = {
  id: 'sql',
  formatter: formatSQL,
}
