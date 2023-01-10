import React, { ChangeEvent, useState } from 'react';
import { LegacyForms, Select, InlineFieldRow, InlineField, Button, Monaco, monacoTypes} from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { FlightSQLDataSource } from '../datasource';
import { FlightSQLDataSourceOptions, SQLQuery } from '../types';
import { useAsync } from 'react-use';
import { QueryEditorRaw } from './RawEditor';
import { LanguageCompletionProvider, getStandardSQLCompletionProvider, } from '@grafana/experimental';
import { formatSQL } from './sqlFormatter'


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

interface CompletionProviderGetterArgs {
  getColumns: React.MutableRefObject<(t: SQLQuery) => Promise<any[]>>;
  getTables: React.MutableRefObject<(d?: string) => Promise<any[]>>;
}

const GetColumns = (datasource: FlightSQLDataSource): AsyncColumnsStateColumn => {
  const result = useAsync(async () => {
    const { columns } = await datasource.getColumns();
    return columns.map((c: any) => ({
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
    const res = await datasource.getTables();
    return res.frames[0].data.values[0].map((t: string) => ({
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

export const getSqlCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider =
  ({ getColumns, getTables }) =>
  (monaco, language) => ({
    ...(language && getStandardSQLCompletionProvider(monaco, {...language, builtinFunctions: [''] })),
    tables: {
      resolve: async () => {
        return await getTables.current();
      },
    },
    columns: {
      // resolve: async (t?: TableIdentifier) => {
        resolve: async (t?: any) => {
        return await getColumns.current(t);
      },
    },
  });


export function QueryEditor(props: QueryEditorProps<FlightSQLDataSource, SQLQuery, FlightSQLDataSourceOptions>) {
  const { onChange, query, datasource } = props;  

  const onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, queryText: event.target.value });
    };

  const onQueryRawSQLChange = (q: SQLQuery) => {
    onChange(q);
    };

  const { queryText } = props.query;
  const { loadingTable, tables, errorTable } = GetTables(datasource);
  const { loadingColumn, columns, errorColumn } = GetColumns(datasource);
  const [table, setTable] = useState<SelectableValue<string>>();
  const [column, setColumn] = useState<SelectableValue<string>>();
  const [builderView, setView] = useState(false);
  const args = {
    getColumns: { current: () =>  datasource.getTables() },
    getTables: { current: () =>  datasource.getColumns() },
  };
  const sqlLanguageDefinition = {
    id: 'sql',
    // id: 'pgsql',
    completionProvider: getSqlCompletionProvider(args),
    formatter: formatSQL,
  };

  return  (
  <>

  {
    builderView ? (
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
    ) : (
<QueryEditorRaw query={query} onChange={onQueryRawSQLChange} editorLanguageDefinition={sqlLanguageDefinition} />
    )
    
  }
   <Button fill="outline" size="sm" onClick={() => setView(!builderView)} >{builderView ? "Raw Editor" :  "Builder View"}</Button>
   </>
  )
}
