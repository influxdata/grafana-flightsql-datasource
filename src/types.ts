import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface SQLQuery extends DataQuery {
  queryText?: string;
}

export const DEFAULT_QUERY: Partial<SQLQuery> = {};

/**
 * These are options configured for each DataSource instance
 */
export interface FlightSQLDataSourceOptions extends DataSourceJsonData {
  host?: string;
  database?: string;
  token?: string;
  secure?: boolean;
}

export type TablesResponse = {
  tables: string[];
};

export type ColumnsResponse = {
  columns: string[];
};
