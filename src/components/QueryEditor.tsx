import React, {useState, useMemo, useCallback, useEffect} from 'react'
import {Button} from '@grafana/ui'
import {QueryEditorProps} from '@grafana/data'
import {FlightSQLDataSource} from '../datasource'
import {FlightSQLDataSourceOptions, SQLQuery} from '../types'

import {QueryEditorRaw} from './QueryEditorRaw'
import {LanguageCompletionProvider, getStandardSQLCompletionProvider} from '@grafana/experimental'
import {formatSQL} from './sqlFormatter'
import {BuilderView} from './BuilderView'

export const COMMON_AGGREGATE_FNS = ['AVG', 'COUNT', 'MAX', 'MIN', 'SUM']
export const getSqlCompletionProvider: (args: any) => LanguageCompletionProvider =
  ({getTables, getColumns, sqlInfo}) =>
  (monaco, language) => {
    return {
      ...(language &&
        getStandardSQLCompletionProvider(monaco, {
          ...language,
          builtinFunctions: sqlInfo.builtinFunctions,
          keywords: sqlInfo.keywords,
        })),
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
      columns: {
        resolve: (t: string) => getColumns(t),
      },
      supportedMacros: () => [],
    }
  }

export function QueryEditor(props: QueryEditorProps<FlightSQLDataSource, SQLQuery, FlightSQLDataSourceOptions>) {
  const {onChange, query, datasource} = props
  const [sqlInfo, setSqlInfo] = useState<any>()
  useEffect(() => {
    ;(async () => {
      const res = await datasource.getSQLInfo()
      const keywords = res?.frames[0].data.values[1][17]
      const numericFunctions = res?.frames[0].data.values[1][18]
      const stringFunctions = res?.frames[0].data.values[1][19]
      const systemFunctions = res?.frames[0].data.values[1][20]
      const sqlDateTimeFunctions = res?.frames[0].data.values[1][21]
      const functions = [...numericFunctions, ...stringFunctions, ...systemFunctions, ...sqlDateTimeFunctions]
      setSqlInfo({keywords: keywords, builtinFunctions: functions})
    })()
  }, [datasource])

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
        sqlInfo,
      }),
    [getTables, getColumns, sqlInfo]
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
          onChange={onChange}
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
