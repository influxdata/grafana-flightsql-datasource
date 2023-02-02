import {DataSourcePluginOptionsEditorProps, PluginType} from '@grafana/data'

import {FlightSQLDataSource} from './datasource'
import {FlightSQLDataSourceOptions, SQLQuery} from './types'

export const mockDatasource = new FlightSQLDataSource({
  id: 1,
  uid: 'influxdata-flightsql-id',
  type: 'influxdata-flightsql-datasource',
  name: 'FlightSQL Data Source',
  readOnly: false,
  jsonData: {},
  access: 'proxy',
  meta: {
    id: 'influxdata-flightsql-datasource',
    module: '',
    name: 'FlightSQL Data Source',
    type: PluginType.datasource,
    alerting: true,
    backend: true,
    baseUrl: 'public/plugins/influxdata-flightsql-datasource',
    info: {
      description: '',
      screenshots: [],
      updated: '',
      version: '',
      logos: {
        small: '',
        large: '',
      },
      author: {
        name: '',
      },
      links: [],
    },
  },
})

export const mockDatasourceOptions: DataSourcePluginOptionsEditorProps<FlightSQLDataSourceOptions> = {
  options: {
    id: 1,
    uid: '1',
    orgId: 1,
    name: 'Timestream',
    typeLogoUrl: '',
    type: '',
    access: '',
    url: '',
    user: '',
    basicAuth: false,
    basicAuthUser: '',
    database: '',
    isDefault: false,
    jsonData: {
      host: '',
      token: '',
      secure: true,
      username: '',
      password: '',
      selectedAuthType: '',
      metadata: [],
    },
    secureJsonFields: {},
    readOnly: false,
    withCredentials: false,
    typeName: '',
  },
  onOptionsChange: jest.fn(),
}

export const mockQuery: SQLQuery = {
  queryText: 'select * from information_schema.tables',
  refId: '',
  format: 'table',
  rawEditor: true,
  table: '',
  columns: [],
  wheres: [],
  orderBy: '',
  groupBy: '',
  limit: '',
}
