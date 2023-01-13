import React, {useEffect, useState} from 'react'
import {css} from '@emotion/css'

import {Select, SegmentSection, InlineLabel, Input} from '@grafana/ui'
import {SelectableValue} from '@grafana/data'
import {GetTables, checkCasing, buildQueryString} from './utils'
import {SelectColumn} from './SelectColumn'

export function BuilderView({query, datasource, onChange}: any) {
  const [columnValues, setColumnValues] = useState([{value: ''}])

  const handleColumnChange = (column: any) => {
    let newColumnValues = [...columnValues]
    newColumnValues[column.index]['value'] = column.value
    setColumnValues(newColumnValues)
  }

  const addColumns = () => {
    setColumnValues([...columnValues, {value: ''}])
  }

  const removeColumns = (i: any) => {
    let newColumnValues = [...columnValues]
    newColumnValues.splice(i, 1)
    setColumnValues(newColumnValues)
  }
  const [groupBy, setGroupBy] = useState('')
  const [where, setWhere] = useState('')
  const [orderBy, setOrderBy] = useState('')
  const [limit, setLimit] = useState('')
  const [table, setTable] = useState<SelectableValue<string>>()
  const [column, setColumn] = useState<SelectableValue<string>>()
  const {loadingTable, tables, errorTable} = GetTables(datasource)
  const [columns, setColumns] = useState()
  const formatCreateLabel = (v: string) => v
  const selectClass = css({
    minWidth: '160px',
  })

  useEffect(() => {
    ;(async () => {
      const res = await datasource.getColumns(table?.value)
      const columns = res.frames[0].schema.fields.map((t: any) => ({
        index: '',
        label: t.name,
        value: t.name,
      }))
      setColumns(columns)
    })()
  }, [table, datasource])

  useEffect(() => {
    if (table && column) {
      const selectColumns = columnValues
        .map((v) => v.value)
        .join(',')
        .replace(/,\s*$/, '')
      const t = checkCasing(table.value || '')
      const queryText = buildQueryString(selectColumns, t, where, orderBy, groupBy, limit)
      onChange({...query, queryText: queryText})
    }
  }, [table, columnValues, groupBy, where, orderBy, limit, column])

  useEffect(() => {
    if (column) {
      handleColumnChange(column)
    }
  }, [column])

  useEffect(() => {
    setColumnValues([{value: ''}])
  }, [table])

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
            width={30}
            placeholder=""
          />
        </SegmentSection>
      </div>
      {columnValues.map((c, index) => {
        return (
          <div className={selectClass} key={index}>
            <SegmentSection label="SELECT" fill={true}>
              <SelectColumn
                columns={columns}
                setColumn={setColumn}
                index={index}
                column={c.value}
                formatCreateLabel={formatCreateLabel}
              />
              <InlineLabel as="button" className="" onClick={addColumns} width="auto">
                +
              </InlineLabel>
              <InlineLabel as="button" className="" width="auto" onClick={removeColumns}>
                -
              </InlineLabel>
            </SegmentSection>
          </div>
        )
      })}
      <div className={selectClass}>
        <SegmentSection label="WHERE" fill={true}>
          <Input
            autoFocus
            type="text"
            spellCheck={false}
            onBlur={() => {}}
            onKeyDown={(e: any) => {
              if (e.key === 'Enter') {
                onChange(where)
              }
            }}
            onChange={(e: any) => {
              setWhere(e.currentTarget.value)
            }}
            value={where}
            width={30}
          />
        </SegmentSection>
      </div>
      <div className={selectClass}>
        <SegmentSection label="GROUP BY" fill={true}>
          <Input
            autoFocus
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
            width={30}
          />
        </SegmentSection>
      </div>
      <div className={selectClass}>
        <SegmentSection label="ORDER BY" fill={true}>
          <Input
            autoFocus
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
            width={30}
          />
        </SegmentSection>
      </div>
      <div className={selectClass}>
        <SegmentSection label="LIMIT" fill={true}>
          <Input
            autoFocus
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
            width={30}
          />
        </SegmentSection>
      </div>
    </>
  )
}
