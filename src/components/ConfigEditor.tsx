import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions } from '../types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onHostChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      host: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onDatabaseChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      database: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      token: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  render() {
    const { options } = this.props;
    const { jsonData } = options;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="Host"
            labelWidth={6}
            inputWidth={20}
            onChange={this.onHostChange}
            value={jsonData.host || ''}
            placeholder=""
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Database"
            labelWidth={6}
            inputWidth={20}
            onChange={this.onDatabaseChange}
            value={jsonData.database || ''}
            placeholder=""
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              value={jsonData.token || ''}
              label="Token"
              placeholder=""
              labelWidth={6}
              inputWidth={20}
              onChange={this.onTokenChange}
              onReset={() => {}}
              isConfigured={false}
            />
          </div>
        </div>
      </div>
    );
  }
}
