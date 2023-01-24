import {useAsync} from 'react-use'
import {FlightSQLDataSource} from '../datasource'
import {SelectableValue} from '@grafana/data'

type AsyncTablesState = {
  loadingTable: boolean
  tables: Array<SelectableValue<string>>
  errorTable: Error | undefined
}

export const getTables = (datasource: FlightSQLDataSource): AsyncTablesState => {
  const result = useAsync(async () => {
    const res = await datasource.getTables()
    return res.frames[0].data.values[2].map((t: string) => ({
      label: t,
      value: t,
    }))
  }, [datasource])

  return {
    loadingTable: result.loading,
    tables: result.value ?? [],
    errorTable: result.error,
  }
}

export const buildQueryString = (
  columns: string,
  table: string,
  whereExp: string,
  orderBy: string | undefined,
  groupBy: string | undefined,
  limit: string | undefined
): string => {
  let queryStr = `SELECT ${columns} FROM ${table}`

  if (whereExp) {
    queryStr = queryStr + ` WHERE ${whereExp}`
  }

  if (groupBy) {
    queryStr = queryStr + ` GROUP BY ${groupBy}`
  }

  if (orderBy) {
    queryStr = queryStr + ` ORDER BY ${orderBy}`
  }

  if (limit) {
    queryStr = queryStr + ` LIMIT ${limit}`
  }

  return queryStr
}

export const checkCasing = (str: string) => {
  const camelCase = /[A-Z]/.test(str)
  if (camelCase) {
    str = `"${str}"`
  }

  return str
}

export const handleColumnChange = (column: any, setColumnValues: any, columnValues: any) => {
  let newColumnValues = [...columnValues]
  newColumnValues[column.index]['value'] = column.value
  setColumnValues(newColumnValues)
}

export const addColumns = (setColumnValues: any, columnValues: any) => {
  setColumnValues([...columnValues, {value: ''}])
}

export const removeColumns = (i: any, setColumnValues: any, columnValues: any) => {
  let newColumnValues = [...columnValues]
  newColumnValues.splice(i, 1)
  setColumnValues(newColumnValues)
}

export const handleWhereChange = (where: any, setWhereValues: any, whereValues: any) => {
  let newWhereValues = [...whereValues]
  newWhereValues[where.index]['value'] = where.value
  setWhereValues(newWhereValues)
}

export const addWheres = (setWhereValues: any, whereValues: any) => {
  setWhereValues([...whereValues, {value: ''}])
}

export const removeWheres = (i: any, setWhereValues: any, whereValues: any) => {
  let newWhereValues = [...whereValues]
  newWhereValues.splice(i, 1)
  setWhereValues(newWhereValues)
}

export const formatColumns = (columnArr: any) => {
  return columnArr
    .map((c: any) => checkCasing(c.value))
    .join(',')
    .replace(/,\s*$/, '')
}

export const formatWheres = (whereValues: any) => {
  return whereValues
    .map((w: any) => w.value)
    .filter(Boolean)
    .join(' and ')
}

export const formatCreateLabel = (v: string) => v
