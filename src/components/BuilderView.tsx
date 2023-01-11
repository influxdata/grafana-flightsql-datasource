import React, {useState, ChangeEvent} from 'react';
import { css } from '@emotion/css';

// import { EditorRows, EditorRow, EditorField, InputGroup } from '@grafana/experimental';
import { Select, LegacyForms, SegmentSection, InlineFieldRow } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';
const { FormField } = LegacyForms;
import { GetTables, GetColumns } from './utils'

export function BuilderView({query, datasource, onChange}: any) {
  const onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, queryText: event.target.value });
    };
  const { queryText } = query;
  const { loadingTable, tables, errorTable } = GetTables(datasource);
  const { loadingColumn, columns, errorColumn } = GetColumns(datasource);
  const [table, setTable] = useState<SelectableValue<string>>();
  const [column, setColumn] = useState<SelectableValue<string>>();
  const formatCreateLabel = (v: string) => v;
  const selectClass = css({
    minWidth: '160px',
  });
  return (
    <>
    {/* <InlineFieldRow> */}
    <FormField
      labelWidth={8}
      inputWidth={30}
      value={queryText || ''}
      onChange={onQueryTextChange}
      label="SQL"
      tooltip="SQL query text"
      />
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
      isLoading={loadingColumn}
      disabled={!!errorColumn}
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
