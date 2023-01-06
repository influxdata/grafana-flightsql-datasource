import React, { ChangeEvent, useState } from 'react';
import { LegacyForms, Select, InlineFieldRow, InlineField } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { FlightSQLDataSource } from '../datasource';
import { FlightSQLDataSourceOptions, SQLQuery } from '../types';
import { useAsync } from 'react-use';

const { FormField } = LegacyForms;

type AsyncTablesState = {
  loadingTable: boolean;
  tables: Array<SelectableValue<string>>;
  errorTable: Error | undefined;
};

type AsyncColumnsStateColumn = {
  loadingColumn: boolean;
  columns: Array<SelectableValue<string>>;
  errorColumn: Error | undefined;
};

const GetColumns = (datasource: FlightSQLDataSource): AsyncColumnsStateColumn => {
  const result = useAsync(async () => {
    const { columns } = await datasource.getColumns();
    return columns.map((c) => ({
      label: c,
      value: c,
    }));
  }, [datasource]);

  return {
    loadingColumn: result.loading,
    columns: result.value ?? [],
    errorColumn: result.error,
  };
}

const GetTables = (datasource: FlightSQLDataSource): AsyncTablesState => {
  const result = useAsync(async () => {
    const { tables } = await datasource.getTables();
    return tables.map((t) => ({
      label: t,
      value: t,
    }));
  }, [datasource]);

  return {
    loadingTable: result.loading,
    tables: result.value ?? [],
    errorTable: result.error,
  };
}

export function QueryEditor(props: QueryEditorProps<FlightSQLDataSource, SQLQuery, FlightSQLDataSourceOptions>) {
  const { onChange, query, datasource } = props;  

  const onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, queryText: event.target.value });
    };
  const { queryText } = props.query;
  const { loadingTable, tables, errorTable } = GetTables(datasource);
  const { loadingColumn, columns, errorColumn } = GetColumns(datasource);
  const [table, setTable] = useState<SelectableValue<string>>();
  const [column, setColumn] = useState<SelectableValue<string>>();

  return  (
  <>
    <FormField
      labelWidth={8}
      inputWidth={30}
      value={queryText || ''}
      onChange={onQueryTextChange}
      label="SQL"
      tooltip="SQL query text"
      />
       <InlineFieldRow>
        <InlineField label="Tables" grow>
          <Select
            options={tables}
            onChange={setTable}
            isLoading={loadingTable}
            disabled={!!errorTable}
            value={table}
          />
        </InlineField>
      </InlineFieldRow>
      <InlineFieldRow>
        <InlineField label="Columns" grow>
          <Select
            options={columns}
            onChange={setColumn}
            isLoading={loadingColumn}
            disabled={!!errorColumn}
            value={column}
          />
        </InlineField>
      </InlineFieldRow>
   </>
  )
}
