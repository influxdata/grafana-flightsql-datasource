import React, { useCallback, useEffect, useRef } from 'react';

import { LanguageDefinition, SQLEditor } from '@grafana/experimental';

// import { SQLQuery } from '../../types';

type Props = {
  query: any;
  onChange: (value: any, processQuery: boolean) => void;
  children?: (props: { formatQuery: () => void }) => React.ReactNode;
  width?: number;
  height?: number;
  editorLanguageDefinition: LanguageDefinition;
};

export function QueryEditorRaw({ children, onChange, query, width, height, editorLanguageDefinition }: Props) {
  // We need to pass query via ref to SQLEditor as onChange is executed via monacoEditor.onDidChangeModelContent callback, not onChange property
  const queryRef = useRef<any>(query);
  useEffect(() => {
    queryRef.current = query;
  }, [query]);

  const onRawQueryChange = useCallback(
    (rawSql: string, processQuery: boolean) => {
      const newQuery = {
        ...queryRef.current,
        rawQuery: true,
        // rawSql,
        queryText: rawSql,
      };
      onChange(newQuery, processQuery);
    },
    [onChange]
  );

  return (
    <SQLEditor
      width={width}
      height={height}
      query={query.rawSql!}
      onChange={onRawQueryChange}
      language={editorLanguageDefinition}
    >
      {children}
    </SQLEditor>
  );
}
