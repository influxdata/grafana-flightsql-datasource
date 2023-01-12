import React, { useState, useMemo, useCallback, useEffect } from 'react';
import { Button} from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { FlightSQLDataSource } from '../datasource';
import { FlightSQLDataSourceOptions, SQLQuery } from '../types';

import { QueryEditorRaw } from './RawEditor';
import { LanguageCompletionProvider, getStandardSQLCompletionProvider, } from '@grafana/experimental';
import { formatSQL } from './sqlFormatter'
import { BuilderView } from './BuilderView';

interface CompletionProviderGetterArgs {
  getColumns: React.MutableRefObject<(t: SQLQuery) => Promise<any[]>>;
  getTables: React.MutableRefObject<(d?: string) => Promise<any[]>>;
}

export const getSqlCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider =
  ({ getColumns, getTables }) =>
  (monaco, language) => ({
    ...(language && getStandardSQLCompletionProvider(monaco, 
      {...language, 
        // todo: make custom for FlightSQL
        // keywords?: string[];
        // builtinFunctions?: string[];
        // logicalOperators?: string[];
        // comparisonOperators?: string[];
        // operators?: string[];
      }
      )),
    tables: {
      resolve: async () => {
        return await getTables.current();
      },
    },
    columns: {
        resolve: async (t?: any) => {
        return await getColumns.current(t);
      },
    },
  });


export function QueryEditor(props: QueryEditorProps<FlightSQLDataSource, SQLQuery, FlightSQLDataSourceOptions>) {
  const { onChange, query, datasource } = props;  

  const onQueryRawSQLChange = (q: SQLQuery) => {
    onChange(q);
    };
 
  const [builderView, setView] = useState(false);
  // todo: fix these 
  const args = {
    getTables: { current: () =>  datasource.getTables() },
    getColumns: { current: (t: any) =>  datasource.getColumns(t) },
  };
  const sqlLanguageDefinition = {
    id: 'sql',
    completionProvider: getSqlCompletionProvider(args),
    formatter: formatSQL,
  };

  return  (
  <>

  {
    builderView ? (
      <BuilderView
      query={props.query}
      datasource={datasource}
      onChange={onChange}
      />
    ) : (
  <QueryEditorRaw  
    query={query} 
    onChange={onQueryRawSQLChange} 
    editorLanguageDefinition={sqlLanguageDefinition} 
  />
    )
    
  }
   <Button fill="outline" size="sm" onClick={() => setView(!builderView)} >{builderView ? "Raw Editor" :  "Builder View"}</Button>
   </>
  )
}
