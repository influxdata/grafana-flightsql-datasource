# Grafana Flight SQL Datasource

⚠️ This library is experimental and under active development. The configuration it provides could change at any time so use at your own risk.

## Running the plugin in Grafana using Docker compose

To run grafana using docker with this plugin loaded in run:

1. To build the frontend and backend into dist:

`yarn build-plugin`

2. To run grafana using docker-compose mounting the dist directory which contains the plugin:

`yarn server`

3. Navigate to [Grafana homepage](http://localhost:3000/) and locate the FlightSQL datasource.

## Running the plugin not using Docker

1. To build the frontend and backend into dist:

`yarn build-plugin`

2. Make a zipfile of the plugin:

`zip grafana-flightsql-datasource-1.0.0.zip dist -r`

3. Unpack the plugin in your chosen directory eg `/Users/myusername/grafana-plugins/` and point your `grafana.ini` paths configuration at it:

```
[paths]
plugins = /Users/myusername/grafana-plugins/
```

OR you will need to set the relevant environment variable:

`GF_PATHS_PLUGINS=/Users/myusername/grafana-plugins/`

4. The plugin is not yet signed so you will need to add it to the grafana configuration of allowed unsigned plugins: `grafana.ini`:

```
[plugins]
allow_loading_unsigned_plugins = grafana-flightsql-datasource
```

OR you will need to set the relevant environment variable:

`GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=grafana-flightsql-datasource`

5. Restart your grafana configuration

6. Navigate to the Grafana homepage and locate the FlightSQL datasource.

Further documentation about grafana configuration can be found [here](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/#plugins)

## Configuring the plugin:

For specific IOX connectivity documentation go to our wiki [here](https://github.com/influxdata/grafana-flightsql-datasource/wiki).

- Host: Provide the host:port of your FlightSQL client.

- Database: Provide your target database name.

- Token: Provide a bearer token for accessing your client.

- Require TLS / SSL: Either enable or disable TLS based on the configuration of your client.

A username/password option is under active development.

## Using the query builder

The default view is a query builder which is in active development:

- Begin by selecting the table from the dropdown.

- This will auto populate your available columns for your select statement. Use the + and - buttons to add or remove additional where statements.

- You can overwrite a dropdown field by typing in your desired value eg `*`.

- The where field is a text entry where you can define any where clauses. Use the + and - buttons to add or remove additional where statements.

- You can switch to a raw sql input by pressing the Edit SQL button this will show you the query you have been building thus far and allow you to enter any query.

- Select the query button to see your results.

- From there you can add to dashboards and create any additional dashboards you like.

## Developing against the plugin

### Frontend setup

1. Install dependencies

   `yarn install`

2. Build plugin in development mode

   `yarn dev`

3. Run the linter prior to merging commits

   `yarn lint:fix`

### Backend setup

1. Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

   ```
   go get -u github.com/grafana/grafana-plugin-sdk-go
   go mod tidy
   ```

2. Build backend plugin binaries for Linux, Windows and Darwin:

   `mage -v`

3. List all available Mage targets for additional commands:

   `mage -l`

4. Run the plugin using docker:

   `yarn server`

Navigate to [Grafana homepage](http://localhost:3000/) and locate the FlightSQL datasource.

## Grafana file structure

plugin.json - defines capabilities of plugin - ie that it's a backend datasource

module.ts - entrypoint for plugin wraps up - DataSourcePlugin

src/components/ConfigEditor.tsx - connection to datasource

src/components/QueryEditor.tsx - query ui feature set

pkg/plugin/main.go - entrypoint into backend

pkg/plugin/datasource.go - backend logic
