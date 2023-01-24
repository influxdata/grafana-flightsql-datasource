# Grafana Flight SQL Datasource

:warning: This library is experimental and under active development. The configuration it provides could change at any time so use at your own risk.

This is a plugin for Grafana that enables queries to Flight SQL APIs.

## Installation

### Download a Release

We haven't released any builds of this plugin yet. You'll need to clone the
repository and build it from source using following the instructions as a guide.

### Building from Source

These instructions assume you have the following tools installed:

- [`yarn`](https://yarnpkg.com)
- [`mage`](https://magefile.org)
- [`grafana version 9.2.x`](https://grafana.com/grafana/download/9.2.0)

#### Docker Compose

Docker Compose is the easiest way to get a Grafana instance running with this
plugin present. A Docker Compose workflow is baked into the `yarn` workflow.

1. Build frontend and backend in `dist/`:

   ```shell
   $ yarn build-plugin
   ```

1. Start Grafana in Docker Compose:

   ```shell
   $ yarn server
   ```

   The artifacts in `dist/` will be mounted and loaded into the installation.

1. Navigate to the [Locally Running Grafana](http://localhost:3000/).
1. Follow the instructions in [Adding a Flight SQL
   Datasource](#adding-a-flight-sql-datasource).

#### Running without Docker

1. Build frontend and backend to produce `influxdata-flightsql-datasource.tar.gz`:

   ```shell
   $ yarn build-plugin
   ```

1. Unpack the archive to your chosen directory (e.g. `grafana-plugins/`).

   ```shell
   $ mkdir -p grafana-plugins/influxdata-flightsql-datasource
   $ tar -xf influxdata-flightsql-datasource.tar.gz -C grafana-plugins/influxdata-flightsql-datasource
   ```

1. Point Grafana to this the plugins directory. You have two options:

   - Edit the `paths.plugins` directive in your `grafana.ini`:

     ```ini
     [paths]
     plugins = grafana-plugins/
     ```

   - **OR** set the relevant environment variable where Grafana is started:
     ```shell
     GF_PATHS_PLUGINS=grafana-plugins/
     ```

1. The plugin is not yet signed so you will need to add it to the Grafana configuration of allowed unsigned plugins.

   - Add the following to your `grafana.ini`:

     ```ini
     [plugins]
     allow_loading_unsigned_plugins = influxdata-flightsql-datasource
     ```

   - **OR** set the relevant environment variable where Grafana is started:
     ```shell
     GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=influxdata-flightsql-datasource
     ```

1. Navigate to the [Locally Running Grafana](http://localhost:3000/).
1. Follow the instructions in [Adding a Flight SQL
   Datasource](#adding-a-flight-sql-datasource).

:book: Further documentation about [Grafana configuration](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/#plugins).

## Usage

### Adding a Flight SQL Datasource

TODO

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
