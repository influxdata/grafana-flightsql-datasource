import React, {useEffect, useState} from 'react'
import {InlineSwitch, FieldSet, InlineField, SecretInput, Input, Select, InlineFieldRow} from '@grafana/ui'
import {DataSourcePluginOptionsEditorProps, SelectableValue} from '@grafana/data'
import {FlightSQLDataSourceOptions, authTypeOptions} from '../types'
import {
  onHostChange,
  onDatabaseChange,
  onTokenChange,
  onSecureChange,
  onUsernameChange,
  onPasswordChange,
  onAuthTypeChange,
} from './utils'

export function ConfigEditor(props: DataSourcePluginOptionsEditorProps<FlightSQLDataSourceOptions>) {
  const {options, onOptionsChange} = props
  const {jsonData} = options
  const [selectedAuthType, setAuthType] = useState<SelectableValue<string>>(authTypeOptions[2])

  useEffect(() => {
    onAuthTypeChange(selectedAuthType, options, onOptionsChange)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedAuthType])

  useEffect(() => {
    const {onOptionsChange, options} = props
    const jsonData = {
      ...options.jsonData,
      selectedAuthType: 'token',
      secure: true,
    }
    onOptionsChange({...options, jsonData})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

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

        <InlineField labelWidth={20} label="Database">
          <Input
            width={40}
            name="database"
            type="text"
            placeholder="dbName"
            onChange={(e) => onDatabaseChange(e, options, onOptionsChange)}
            value={jsonData.database || ''}
          ></Input>
        </InlineField>
        <InlineField labelWidth={20} label="Auth Type">
          <Select
            options={authTypeOptions}
            onChange={setAuthType}
            value={jsonData.selectedAuthType || ''}
            allowCustomValue={true}
            autoFocus={true}
            // formatCreateLabel={}
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
          <InlineFieldRow>
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
    </div>
  )
}
