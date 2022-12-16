import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  queryText?: string;
}

export const DEFAULT_QUERY: Partial<MyQuery> = {};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  host?: string;
  database?: string;
  token?: string;
}
