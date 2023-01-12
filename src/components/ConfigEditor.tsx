import React, {ChangeEvent, PureComponent} from 'react'
import {InlineSwitch, FieldSet, InlineField, SecretInput, Input} from '@grafana/ui'
import {DataSourcePluginOptionsEditorProps} from '@grafana/data'
import {FlightSQLDataSourceOptions} from '../types'

interface Props extends DataSourcePluginOptionsEditorProps<FlightSQLDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onHostChange = (event: ChangeEvent<HTMLInputElement>) => {
    const {onOptionsChange, options} = this.props
    const jsonData = {
      ...options.jsonData,
      host: event.target.value,
    }
    onOptionsChange({...options, jsonData})
  }

  onDatabaseChange = (event: ChangeEvent<HTMLInputElement>) => {
    const {onOptionsChange, options} = this.props
    const jsonData = {
      ...options.jsonData,
      database: event.target.value,
    }
    onOptionsChange({...options, jsonData})
  }

  onTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    const {onOptionsChange, options} = this.props
    const jsonData = {
      ...options.jsonData,
      token: event.target.value,
    }
    onOptionsChange({...options, jsonData})
  }

  onSecureChange = () => {
    const {onOptionsChange, options} = this.props
    const jsonData = {
      ...options.jsonData,
      secure: !options.jsonData.secure,
    }
    onOptionsChange({...options, jsonData})
  }

  render() {
    const {options} = this.props
    const {jsonData} = options
    return (
      <div>
        <FieldSet label="FlightSQL Connection" width={400}>
          <InlineField labelWidth={20} label="Host">
            <Input
              width={40}
              name="host"
              type="text"
              value={jsonData.host || ''}
              placeholder="localhost:1234"
              onChange={this.onHostChange}
            ></Input>
          </InlineField>

          <InlineField labelWidth={20} label="Database">
            <Input
              width={40}
              name="database"
              type="text"
              placeholder="dbName"
              onChange={this.onDatabaseChange}
              value={jsonData.database || ''}
            ></Input>
          </InlineField>
          <InlineField labelWidth={20} label="Token">
            <SecretInput
              width={40}
              name="token"
              type="text"
              value={jsonData.token || ''}
              placeholder="****************"
              onChange={this.onTokenChange}
              onReset={() => {}}
              isConfigured={false}
            ></SecretInput>
          </InlineField>

          <InlineField labelWidth={20} label="Require TLS / SSL">
            <InlineSwitch
              label=""
              value={jsonData.secure}
              onChange={this.onSecureChange}
              showLabel={false}
              disabled={false}
            />
          </InlineField>
        </FieldSet>
      </div>
    )
  }
}
