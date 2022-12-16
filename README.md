# Grafana Flight SQL Datasource

The inital template for this datasoure plugin was built using the following command, choosing backend datasource.

`npx @grafana/create-plugin`

## Getting started

To run grafana using docker with this plugin loaded in run:

`yarn server`

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

Navigate to [Grafana homepage](http://localhost:3000/) and locate your datasource.