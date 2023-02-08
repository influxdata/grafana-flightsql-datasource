import React, {useEffect, useState} from 'react'
import {css} from '@emotion/css'

import {Select, SegmentSection, InlineLabel, Input} from '@grafana/ui'
import {SelectableValue} from '@grafana/data'
import {
  GetTables,
  checkCasing,
  buildQueryString,
  handleColumnChange,
  addColumns,
  removeColumns,
  handleWhereChange,
  addWheres,
  removeWheres,
  formatColumns,
  formatWheres,
  formatCreateLabel,
} from './utils'
import {SelectColumn} from './SelectColumn'
import {WhereExp} from './WhereExp'

export function BuilderView({query, datasource, onChange}: any) {
  const [columnValues, setColumnValues] = useState([{value: ''}])
  const [whereValues, setWhereValues] = useState([{value: ''}])
  const [groupBy, setGroupBy] = useState('')
  const [where, setWhere] = useState('')
  const [orderBy, setOrderBy] = useState('')
  const [limit, setLimit] = useState('')
  const [columns, setColumns] = useState()
  const [table, setTable] = useState<SelectableValue<string>>()
  const [column, setColumn] = useState<SelectableValue<string>>()

  const {loadingTable, tables, errorTable} = GetTables(datasource)

  useEffect(() => {
    ;(async () => {
      let res
      if (table?.value) {
        res = table?.value && (await datasource.getColumns(table?.value))
      }
      const columns = res?.frames[0].schema.fields.map((t: any) => ({
        index: '',
        label: t.name,
        value: t.name,
      }))
      setColumns(columns)
    })()
  }, [table, datasource])

  useEffect(() => {
    if (table && column) {
      const selectColumns = formatColumns(columnValues)
      const casedTable = checkCasing(table.value || '')
      const whereExps = formatWheres(whereValues)
      const queryText = buildQueryString(selectColumns, casedTable, whereExps, orderBy, groupBy, limit)
      onChange({...query, queryText: queryText})
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [table, columnValues, groupBy, whereValues, orderBy, limit, column])

  useEffect(() => {
    if (column) {
      handleColumnChange(column, setColumnValues, columnValues)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [column])

  useEffect(() => {
    if (where) {
      handleWhereChange(where, setWhereValues, whereValues)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [where])

  useEffect(() => {
    setColumnValues([{value: ''}])
  }, [table])

  const selectClass = css({
    minWidth: '160px',
  })
  return (
    <>
      <div className={selectClass}>
        <SegmentSection label="FROM" fill={true}>
          <Select
            options={tables}
            onChange={setTable}
            isLoading={loadingTable}
            disabled={!!errorTable}
            value={table}
            allowCustomValue={true}
            autoFocus={true}
            formatCreateLabel={formatCreateLabel}
            width={15}
            placeholder="table"
          />
        </SegmentSection>
      </div>
      <div className={selectClass}>
        <SegmentSection label="SELECT" fill={true}>
          {columnValues.map((_, index) => (
            <>
              <SelectColumn
                columns={columns}
                setColumn={setColumn}
                index={index}
                columnValues={columnValues}
                formatCreateLabel={formatCreateLabel}
              />
              {index + 1 >= columnValues.length && (
                <InlineLabel
                  as="button"
                  className=""
                  onClick={() => addColumns(setColumnValues, columnValues)}
                  width="auto"
                  style={{marginLeft: '5px'}}
                >
                  +
                </InlineLabel>
              )}

              {index > 0 && (
                <InlineLabel
                  as="button"
                  className=""
                  width="auto"
                  onClick={() => removeColumns(index, setColumnValues, columnValues)}
                >
                  -
                </InlineLabel>
              )}
            </>
          ))}
        </SegmentSection>
      </div>
      <div className={selectClass}>
        <SegmentSection label="WHERE" fill={true}>
          {whereValues.map((_, index) => (
            <>
              <WhereExp index={index} setWhere={setWhere} whereValues={whereValues} />
              {index + 1 >= whereValues.length && (
                <InlineLabel
                  as="button"
                  className=""
                  onClick={() => addWheres(setWhereValues, whereValues)}
                  width="auto"
                  style={{marginLeft: '5px'}}
                >
                  +
                </InlineLabel>
              )}
              {index > 0 && (
                <InlineLabel
                  as="button"
                  className=""
                  width="auto"
                  onClick={() => removeWheres(index, setWhereValues, whereValues)}
                >
                  -
                </InlineLabel>
              )}
            </>
          ))}
        </SegmentSection>
      </div>
      <div className={selectClass}>
        <SegmentSection label="GROUP BY" fill={true}>
          <Input
            type="text"
            spellCheck={false}
            onBlur={() => {}}
            onKeyDown={(e: any) => {
              if (e.key === 'Enter') {
                onChange(groupBy)
              }
            }}
            onChange={(e: any) => {
              setGroupBy(e.currentTarget.value)
            }}
            value={groupBy}
            width={15}
            placeholder="(optional)"
          />
        </SegmentSection>
      </div>
      <div className={selectClass}>
        <SegmentSection label="ORDER BY" fill={true}>
          <Input
            type="text"
            spellCheck={false}
            onBlur={() => {}}
            onKeyDown={(e: any) => {
              if (e.key === 'Enter') {
                onChange(orderBy)
              }
            }}
            onChange={(e: any) => {
              setOrderBy(e.currentTarget.value)
            }}
            value={orderBy}
            width={15}
            placeholder="(optional)"
          />
        </SegmentSection>
      </div>
      <div className={selectClass}>
        <SegmentSection label="LIMIT" fill={true}>
          <Input
            autoFocus={false}
            type="text"
            spellCheck={false}
            onBlur={() => {}}
            onKeyDown={(e: any) => {
              if (e.key === 'Enter') {
                onChange(limit)
              }
            }}
            onChange={(e: any) => {
              setLimit(e.currentTarget.value)
            }}
            value={limit}
            width={15}
            placeholder="(optional)"
          />
        </SegmentSection>
      </div>
    </>
  )
}
