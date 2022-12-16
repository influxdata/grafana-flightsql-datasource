# Grafana Flight SQL Datasource

The inital template for this datasoure plugin was built using the following command, choosing backend datasource.

`npx @grafana/create-plugin`

## Getting started

Create a directory named grafana-plugins and place this datasource under it.

If running grafana via docker the following command will build grafana with all plugins in that directory:

```
docker run -d -p 3000:3000 -v "$(pwd)"/grafana-plugins:/var/lib/grafana/plugins --name=grafana grafana/grafana:7.0.0
```

### Frontend setup

1. Install dependencies

   ```bash
   yarn install
   ```

2. Build plugin in development mode or run in watch mode

   ```bash
   yarn dev
   ```

3. Run the linter

   ```bash
   yarn lint:fix
   ```

### Backend setup

1. Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

   ```bash
   go get -u github.com/grafana/grafana-plugin-sdk-go
   go mod tidy
   ```

2. Build backend plugin binaries for Linux, Windows and Darwin:

   ```bash
   mage -v
   ```

3. List all available Mage targets for additional commands:

   ```bash
   mage -l
   ```

Navigate to [Grafana homepage](http://localhost:3000/) and locate your datasource.