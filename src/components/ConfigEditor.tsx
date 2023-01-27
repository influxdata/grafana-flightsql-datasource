import React, {useEffect, useState} from 'react'
import {InlineSwitch, FieldSet, InlineField, SecretInput, Input, Select, InlineFieldRow, InlineLabel} from '@grafana/ui'
import {DataSourcePluginOptionsEditorProps, SelectableValue} from '@grafana/data'
import {FlightSQLDataSourceOptions, authTypeOptions} from '../types'
import {
  onHostChange,
  onTokenChange,
  onSecureChange,
  onUsernameChange,
  onPasswordChange,
  onAuthTypeChange,
  onKeyChange,
  onValueChange,
  addMetaData,
  removeMetaData,
} from './utils'

export function ConfigEditor(props: DataSourcePluginOptionsEditorProps<FlightSQLDataSourceOptions>) {
  const {options, onOptionsChange} = props
  const {jsonData} = options
  const [selectedAuthType, setAuthType] = useState<SelectableValue<string>>(
    {value: jsonData.selectedAuthType, label: jsonData.selectedAuthType} || authTypeOptions[2]
  )
  const existingMetastate = jsonData?.metadata?.map((m: any) => ({key: Object.keys(m)[0], value: Object.values(m)[0]}))
  const [metaDataArr, setMetaData] = useState(existingMetastate || [{key: '', value: ''}])
  useEffect(() => {
    onAuthTypeChange(selectedAuthType, options, onOptionsChange)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedAuthType])

  useEffect(() => {
    const {onOptionsChange, options} = props
    const jsonData = {
      ...options.jsonData,
      secure: true,
    }
    onOptionsChange({...options, jsonData})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    const {onOptionsChange, options} = props
    const mapData = metaDataArr?.map((m: any) => ({[m.key]: m.value}))
    const jsonData = {
      ...options.jsonData,
      metadata: mapData,
    }
    onOptionsChange({...options, jsonData})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [metaDataArr])

  return (
    <div>
      <FieldSet label="FlightSQL Connection" width={400}>
        <InlineField labelWidth={20} label="Host:Port">
          <Input
            width={40}
            name="host"
            type="text"
            value={jsonData.host || ''}
            placeholder="localhost:1234"
            onChange={(e) => onHostChange(e, options, onOptionsChange)}
          ></Input>
        </InlineField>
        <InlineField labelWidth={20} label="Auth Type">
          <Select
            options={authTypeOptions}
            onChange={setAuthType}
            value={selectedAuthType || ''}
            allowCustomValue={true}
            width={40}
            placeholder="token"
          />
        </InlineField>
        {selectedAuthType?.label === 'token' && (
          <InlineField labelWidth={20} label="Token">
            <SecretInput
              width={40}
              name="token"
              type="text"
              value={jsonData.token || ''}
              placeholder="****************"
              onChange={(e) => onTokenChange(e, options, onOptionsChange)}
              onReset={() => onTokenChange(null, options, onOptionsChange)}
              isConfigured={false}
            ></SecretInput>
          </InlineField>
        )}
        {selectedAuthType?.label === 'username/password' && (
          <InlineFieldRow style={{flexFlow: 'row'}}>
            <InlineField labelWidth={20} label="Username">
              <Input
                width={40}
                name="username"
                type="text"
                placeholder="username"
                onChange={(e) => onUsernameChange(e, options, onOptionsChange)}
                value={jsonData.username || ''}
              ></Input>
            </InlineField>
            <InlineField labelWidth={20} label="Password">
              <SecretInput
                width={40}
                name="password"
                type="text"
                value={jsonData.password || ''}
                placeholder="****************"
                onChange={(e) => onPasswordChange(e, options, onOptionsChange)}
                onReset={() => onPasswordChange(null, options, onOptionsChange)}
                isConfigured={false}
              ></SecretInput>
            </InlineField>
          </InlineFieldRow>
        )}

        <InlineField labelWidth={20} label="Require TLS / SSL">
          <InlineSwitch
            label=""
            value={jsonData.secure}
            onChange={() => onSecureChange(options, onOptionsChange)}
            showLabel={false}
            disabled={false}
          />
        </InlineField>
      </FieldSet>
      <FieldSet label="MetaData" width={400}>
        {metaDataArr?.map((_: any, i: any) => (
          <InlineFieldRow key={i} style={{flexFlow: 'row'}}>
            <InlineField labelWidth={20} label="Key">
              <Input
                key={i}
                width={40}
                name="key"
                type="text"
                value={metaDataArr[i].key || ''}
                placeholder="key"
                onChange={(e) => onKeyChange(e, options, onOptionsChange, metaDataArr, i, setMetaData)}
              ></Input>
            </InlineField>
            <InlineField labelWidth={20} label="Value">
              <Input
                key={i}
                width={40}
                name="value"
                type="text"
                value={metaDataArr[i].value || ''}
                placeholder="value"
                onChange={(e) => onValueChange(e, options, onOptionsChange, metaDataArr, i, setMetaData)}
              ></Input>
            </InlineField>
            {i + 1 >= metaDataArr.length && (
              <InlineLabel as="button" className="" onClick={() => addMetaData(setMetaData, metaDataArr)} width="auto">
                +
              </InlineLabel>
            )}
            {i > 0 && (
              <InlineLabel
                as="button"
                className=""
                width="auto"
                onClick={() => removeMetaData(i, setMetaData, metaDataArr)}
              >
                -
              </InlineLabel>
            )}
          </InlineFieldRow>
        ))}
      </FieldSet>
    </div>
  )
}
