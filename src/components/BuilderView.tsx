import React, {useEffect, useState} from 'react';
// import { useAsync } from 'react-use';
import { css } from '@emotion/css';

import { Select, SegmentSection } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
// const { FormField } = LegacyForms;
import { GetTables } from './utils'

export function BuilderView({query, datasource, onChange}: any) {
  // const onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   onChange({ ...query, queryText: event.target.value });
  //   };
  // const { queryText } = query;
  const [table, setTable] = useState<SelectableValue<string>>();
  const [column, setColumn] = useState<SelectableValue<string>>();
  const { loadingTable, tables, errorTable } = GetTables(datasource);
  // const { loadingColumn, columns, errorColumn } = GetColumns(datasource, table?.value! ||
  //   '');
  const [columns, setColumns] = useState()
  const formatCreateLabel = (v: string) => v;
  const selectClass = css({
    minWidth: '160px',
  });
  useEffect(() => {
    (async () => {
      const res = await datasource.getColumns(table?.value);
      console.log("res", res.frames[0])
      const columns = res.frames[0].data.values[0].map((t: string) => ({
        label: t,
        value: t,
      }))
      setColumns(columns);
    })();

  }, [table])
  return (
    <>
    {/* <InlineFieldRow> */}
    {/* <FormField
      labelWidth={8}
      inputWidth={30}
      value={queryText || ''}
      onChange={onQueryTextChange}
      label="SQL"
      tooltip="SQL query text"
      /> */}
      {/* </InlineFieldRow>
      <InlineFieldRow> */}
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
        width={60}
        placeholder=""

      />
    </SegmentSection>
    </div>
    {/* </InlineFieldRow>
    <InlineFieldRow> */}
    <div className={selectClass}>
    <SegmentSection label="SELECT" fill={true}>
    <Select
      options={columns}
      onChange={setColumn}
      // isLoading={loadingColumn}
      // disabled={!!errorColumn}
      value={column}
      allowCustomValue={true}
      autoFocus={true}
      formatCreateLabel={formatCreateLabel}
      width={60}
      placeholder=""
    />
    </SegmentSection>
    </div>
    {/* </InlineFieldRow>
    <InlineFieldRow> */}
    <div className={selectClass}>
    <SegmentSection label="WHERE" fill={true}>
    <Select
     options={[]}
     onChange={() => {}}
     isLoading={false}
     disabled={false}
     value={""}
     allowCustomValue={true}
     autoFocus={true}
     formatCreateLabel={formatCreateLabel}
     width={60}
     placeholder=""
    />
    </SegmentSection>
    </div>
    {/* </InlineFieldRow>
    <InlineFieldRow> */}
    <div className={selectClass}>
    <SegmentSection label="GROUP BY" fill={true}>
    <Select
      options={[]}
      onChange={() => {}}
      isLoading={false}
      disabled={false}
      value={""}
      allowCustomValue={true}
      autoFocus={true}
      formatCreateLabel={formatCreateLabel}
      width={60}
      placeholder=""
    />
    </SegmentSection>
    </div>
    {/* </InlineFieldRow> */}
    </>
  );

}
