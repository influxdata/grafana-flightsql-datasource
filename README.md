:warning: This plugin is not under active development but will be kept up to date with security and bug fixes. Feature PRs should be directed towards the [InfluxDB Grafana](https://github.com/grafana/grafana/tree/main/pkg/tsdb/influxdb) plugin which is now the recommended way to use SQL with InfluxDB 3.0.

# Grafana Flight SQL Datasource

This is a plugin for Grafana that enables queries to Flight SQL APIs.

## Requirements

The plugin requires the user to run Grafana >=9.2.5.

## Install via Grafana Cloud
Because Grafana Flight SQL Plugin is a community published plugin, you can install it to your existing Grafana Cloud instance.

- In your cloud instance navigate to connections and search for FlightSQL, follow the link to install.
- Or directly go to this [link](https://grafana.com/grafana/plugins/influxdata-flightsql-datasource/?tab=installation) to install the plugin into your cloud instance.
- Now in your instance you can navigate back to connections and search for FlightSQL which will be installed and allow you to create a FlightSQL datasource.

## Local installation 

Download the [latest Grafana Flight SQL plugin](https://github.com/influxdata/grafana-flightsql-datasource/releases).

Create a directory for grafana to access your custom-plugins eg custom/plugins/directory.

The following shell script downloads and extracts the latest Flight SQL plugin source code into the the current working directory. Run the following inside your grafana plugin directory:

`sh download-grafana-flightsql-plugin.sh`

### Install local grafana

1. Point your local instance of Grafana to the plugins directory. You have two options:

   - Edit the `paths.plugins` directive in your `grafana.ini`:

     ```ini
     [paths]
     plugins = custom/plugins/directory/
     ```

   - **OR** set the relevant environment variable where Grafana is started:
     ```shell
     GF_PATHS_PLUGINS=custom/plugins/directory/
     ```

1. Navigate to the [Locally Running Grafana](http://localhost:3000/).
1. Follow the instructions in [Adding a Flight SQL
   Datasource](#adding-a-flight-sql-datasource).
   
## Install with Docker Run

```
docker run \
  --volume $PWD/influxdata-flightsql-datasource:/custom/plugins/directory/influxdata-flightsql-datasource \
  --publish 3000:3000 \
  --name grafana \
  grafana/grafana:latest
```

## Install with Docker-Compose
```
version: '3'
services:
  grafana:
    image: grafana/grafana:latest
    ports:
      - 3000:3000
    volumes: 
      - ./influxdata-flightsql-datasource:/custom/plugins/directory/influxdata-flightsql-datasource
    restart: always
```

## Note if using versions <=0.1.9 of the plugin

1. The plugin was only signed at version 1.0.0 you will need to add the plugin to your grafana configuration under allow_loading_unsigned_plugins.

   - Add the following to your `grafana.ini`:

     ```ini
     [plugins]
     allow_loading_unsigned_plugins = influxdata-flightsql-datasource
     ```

   - **OR** set the relevant environment variable where Grafana is started:
     ```shell
     GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=influxdata-flightsql-datasource
     ```

## Usage

### Adding a Flight SQL Datasource

1. Open the side menu by clicking the Grafana icon in the top header.
1. In the side menu under the Dashboards link you should find a link named Data Sources.
1. Click the + Add data source button in the top header.
1. Select FlightSQL from the Type dropdown.

### Configuring the Plugin

- **Host:** Provide the host:port of your Flight SQL client.
- **AuthType** Select between none, username/password and token.
- **Token:** If auth type is token provide a bearer token for accessing your client.
- **Username/Password** iF auth type is username and password provide a username and password.
- **Require TLS/SSL:** Either enable or disable TLS based on the configuration of your client.

- **MetaData** Provide optional key, value pairs that you need sent to your Flight SQL client.

Vendor-specific connectivity documentation can be [found in the wiki](https://github.com/influxdata/grafana-flightsql-datasource/wiki).

### Using the Query Builder

The default view is a query builder which is in active development:

- Begin by selecting the table from the dropdown.
- This will auto populate your available columns for your select statement. Use the **+** and **-** buttons to add or remove additional where statements.
- You can overwrite a dropdown field by typing in your desired value (e.g. `*`).
- The where field is a text entry where you can define any where clauses. Use the + and - buttons to add or remove additional where statements.
- You can switch to a raw SQL input by pressing the "Edit SQL" button. This will show you the query you have been building thus far and allow you to enter any query.
- Press the "Run query" button to see your results.
- From there you can add to dashboards and create any additional dashboards you like.

## Development

See [DEVELOPMENT.md](DEVELOPMENT.md).
