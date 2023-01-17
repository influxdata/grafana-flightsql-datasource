import React, {ChangeEvent, useEffect, useState} from 'react'
import {InlineSwitch, FieldSet, InlineField, SecretInput, Input, Select} from '@grafana/ui'
import {DataSourcePluginOptionsEditorProps, SelectableValue} from '@grafana/data'
import {FlightSQLDataSourceOptions} from '../types'

export function ConfigEditor(props: DataSourcePluginOptionsEditorProps<FlightSQLDataSourceOptions>) {
  const onHostChange = (event: ChangeEvent<HTMLInputElement>) => {
    const {onOptionsChange, options} = props
    const jsonData = {
      ...options.jsonData,
      host: event.target.value,
    }
    onOptionsChange({...options, jsonData})
  }

  const onDatabaseChange = (event: ChangeEvent<HTMLInputElement>) => {
    const {onOptionsChange, options} = props
    const jsonData = {
      ...options.jsonData,
      database: event.target.value,
    }
    onOptionsChange({...options, jsonData})
  }

  const onTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    const {onOptionsChange, options} = props
    const jsonData = {
      ...options.jsonData,
      token: event.target.value,
    }
    onOptionsChange({...options, jsonData})
  }

  const onSecureChange = () => {
    const {onOptionsChange, options} = props
    const jsonData = {
      ...options.jsonData,
      secure: !options.jsonData.secure,
    }
    onOptionsChange({...options, jsonData})
  }

  const onUsernameChange = (event: ChangeEvent<HTMLInputElement>) => {
    const {onOptionsChange, options} = props
    const jsonData = {
      ...options.jsonData,
      username: event.target.value,
    }
    onOptionsChange({...options, jsonData})
  }

  const onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    const {onOptionsChange, options} = props
    const jsonData = {
      ...options.jsonData,
      password: event.target.value,
    }
    onOptionsChange({...options, jsonData})
  }
  const authTypeOptions = [
    {key: 0, label: 'none', title: 'none'},
    {key: 1, label: 'username/password', title: 'username/password'},
    {key: 2, label: 'token', title: 'token'},
  ]
  const [selectedAuthType, setAuthType] = useState<SelectableValue<string>>(authTypeOptions[2])
  const onAuthTypeChange = () => {
    const {onOptionsChange, options} = props
    const jsonData = {
      ...options.jsonData,
      selectedAuthType: selectedAuthType?.value,
    }
    onOptionsChange({...options, jsonData})
  }

  useEffect(() => {
    onAuthTypeChange()
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

  const {options} = props
  const {jsonData} = options

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
            onChange={onHostChange}
          ></Input>
        </InlineField>

        <InlineField labelWidth={20} label="Database">
          <Input
            width={40}
            name="database"
            type="text"
            placeholder="dbName"
            onChange={onDatabaseChange}
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
              onChange={onTokenChange}
              onReset={() => {}}
              isConfigured={false}
            ></SecretInput>
          </InlineField>
        )}
        {selectedAuthType?.label === 'username/password' && (
          <>
            <InlineField labelWidth={20} label="Username">
              <Input
                width={40}
                name="username"
                type="text"
                placeholder="username"
                onChange={onUsernameChange}
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
                onChange={onPasswordChange}
                onReset={() => {}}
                isConfigured={false}
              ></SecretInput>
            </InlineField>
          </>
        )}

        <InlineField labelWidth={20} label="Require TLS / SSL">
          <InlineSwitch label="" value={jsonData.secure} onChange={onSecureChange} showLabel={false} disabled={false} />
        </InlineField>
      </FieldSet>
    </div>
  )
}
