// import { css } from '@emotion/css';
import React, { useCallback, useEffect, useRef } from 'react';

import { LanguageDefinition, SQLEditor } from '@grafana/experimental';
// import { useStyles2, useTheme2 } from '@grafana/ui';
// import { GrafanaTheme2 } from '@grafana/data';

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
  // const theme = useTheme2();
  // const styles = useStyles2(getStyles);
  const queryRef = useRef<any>(query);
  useEffect(() => {
    queryRef.current = query;
  }, [query]);

  const onRawQueryChange = useCallback(
    (rawSql: string, processQuery: boolean) => {
      const newQuery = {
        ...queryRef.current,
        rawQuery: true,
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
      key="1"
    >
      {children}
    </SQLEditor>
  );
}

// function getStyles(theme: GrafanaTheme2) {
//   return {
//     modal: css`
//       width: 95vw;
//       height: 95vh;
//     `,
//     modalContent: css`
//       height: 100%;
//       padding-top: 0;
//     `,
//   };
// }