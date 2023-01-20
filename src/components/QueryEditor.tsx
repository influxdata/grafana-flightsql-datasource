import React, {useState, useMemo, useCallback, useEffect} from 'react'
import {Button, Modal} from '@grafana/ui'
import {QueryEditorProps} from '@grafana/data'
import {FlightSQLDataSource} from '../datasource'
import {FlightSQLDataSourceOptions, SQLQuery} from '../types'

import {QueryEditorRaw} from './QueryEditorRaw'
import {LanguageCompletionProvider, getStandardSQLCompletionProvider, MacroType} from '@grafana/experimental'
import {formatSQL} from './sqlFormatter'
import {BuilderView} from './BuilderView'

export const getSqlCompletionProvider: (args: any) => LanguageCompletionProvider =
  ({getTables, getColumns, sqlInfo, macros}) =>
  (monaco, language) => {
    return {
      ...(language &&
        getStandardSQLCompletionProvider(monaco, {
          ...language,
          builtinFunctions: sqlInfo.builtinFunctions,
          keywords: sqlInfo.keywords,
        })),
      triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
      tables: {
        resolve: () => {
          return getTables()
        },
      },
      columns: {
        resolve: (t: string) => getColumns(t),
      },
      supportedMacros: () => macros,
    }
  }

export function QueryEditor(props: QueryEditorProps<FlightSQLDataSource, SQLQuery, FlightSQLDataSourceOptions>) {
  const {onChange, query, datasource} = props
  const [isExpanded, setIsExpanded] = useState(false)
  const [sqlInfo, setSqlInfo] = useState<any>()
  const [macros, setMacros] = useState<any>()

  useCallback(() => {
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

  useEffect(() => {
    ;(async () => {
      const res = await datasource.getMacros()
      const prefix = `$__`
      const macroArr = res?.macros.map((m: any) => prefix.concat(m))
      const macros = macroArr.map((m: any) => ({text: m, name: m, id: m, type: MacroType.Value, args: []}))
      setMacros(macros)
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
        macros,
      }),
    [getTables, getColumns, sqlInfo, macros]
  )

  const sqlLanguageDefinition = {
    id: 'sql',
    formatter: formatSQL,
  }

  return (
    <>
      {isExpanded && (
        <Modal
          title="Warning"
          closeOnBackdropClick={false}
          closeOnEscape={false}
          isOpen={isExpanded}
          onDismiss={() => {
            setIsExpanded(false)
          }}
        >
          {builderView
            ? 'By switching to the raw sql editor if you click to come back to the builder view you will need to refill your query.'
            : 'By switching to the builder view you will not bring your current raw query over to the builder editor, you will have to fill it out again.'}
          <Modal.ButtonRow>
            <Button fill="solid" size="md" variant="secondary" onClick={() => setIsExpanded(!isExpanded)}>
              Back
            </Button>
            <Button
              fill="solid"
              size="md"
              variant="destructive"
              onClick={() => {
                setIsExpanded(!isExpanded)
                setView(!builderView)
              }}
            >
              Switch
            </Button>
          </Modal.ButtonRow>
        </Modal>
      )}
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
      <Button fill="outline" size="sm" onClick={() => setIsExpanded(!isExpanded)}>
        {builderView ? 'Edit SQL' : 'Builder View'}
      </Button>
    </>
  )
}
