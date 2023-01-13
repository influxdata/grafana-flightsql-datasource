import React, {useState, useMemo, useCallback} from 'react'
import {Button} from '@grafana/ui'
import {QueryEditorProps} from '@grafana/data'
import {FlightSQLDataSource} from '../datasource'
import {FlightSQLDataSourceOptions, SQLQuery} from '../types'

import {QueryEditorRaw} from './QueryEditorRaw'
import {LanguageCompletionProvider, getStandardSQLCompletionProvider} from '@grafana/experimental'
import {formatSQL} from './sqlFormatter'
import {BuilderView} from './BuilderView'

// interface CompletionProviderGetterArgs {
//   // getSchemas: () => Promise<any[]>;
//   getTables: () => Promise<any[]>;
//   getColumns: (table?: string) => Promise<any[]>;
// }

/// todo: make custom for FlightSQL
// keywords?: string[];
// builtinFunctions?: string[];
// logicalOperators?: string[];
// comparisonOperators?: string[];
// operators?: string[];

export const COMMON_AGGREGATE_FNS = ['AVG', 'COUNT', 'MAX', 'MIN', 'SUM']
export const getSqlCompletionProvider: (args: any) => LanguageCompletionProvider =
  ({getTables, getColumns}) =>
  (monaco, language) => {
    return {
      ...(language && getStandardSQLCompletionProvider(monaco, {...language, builtinFunctions: COMMON_AGGREGATE_FNS})),
      triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
      // would this be useful
      // schemas: {
      //   resolve: getSchemas,
      // },
      tables: {
        resolve: () => {
          return getTables()
        },
      },
      // how would this be implemented in this view
      // as the table selection
      // comes after the column selection
      columns: {
        resolve: (t: string) => getColumns(t),
      },
      supportedMacros: () => [],
    }
  }

export function QueryEditor(props: QueryEditorProps<FlightSQLDataSource, SQLQuery, FlightSQLDataSourceOptions>) {
  const {onChange, query, datasource} = props

  const onQueryRawSQLChange = (q: SQLQuery) => {
    onChange(q)
  }

  const [builderView, setView] = useState(true)

  const getTables = useCallback(async () => {
    const res = await datasource.getTables()
    return res.frames[0].data.values[2].map((t: string) => ({
      name: t,
    }))
  }, [datasource])

  const getColumns = useCallback(
    async (table: any) => {
      let res
      if (table?.value) {
        res = await datasource.getColumns(table?.value)
      }
      return res?.frames[0].schema.fields
    },
    [datasource]
  )

  const completionProvider = useMemo(
    () =>
      getSqlCompletionProvider({
        getTables,
        getColumns,
      }),
    [getTables, getColumns]
  )

  const sqlLanguageDefinition = {
    id: 'sql',
    formatter: formatSQL,
  }

  return (
    <>
      {builderView ? (
        <BuilderView query={props.query} datasource={datasource} onChange={onChange} />
      ) : (
        <QueryEditorRaw
          query={query}
          onChange={onQueryRawSQLChange}
          editorLanguageDefinition={{
            ...sqlLanguageDefinition,
            completionProvider,
          }}
        />
      )}
      <Button fill="outline" size="sm" onClick={() => setView(!builderView)}>
        {builderView ? 'Edit SQL' : 'Builder View'}
      </Button>
    </>
  )
}
