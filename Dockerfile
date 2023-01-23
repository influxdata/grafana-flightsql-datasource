FROM grafana/grafana:9.2.5

ENV GF_AUTH_DISABLE_LOGIN_FORM "true"

ENV GF_AUTH_ANONYMOUS_ENABLED "true"

ENV GF_AUTH_ANONYMOUS_ORG_ROLE "Admin"

ENV GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=grafana-flightsql-datasource

ADD ./dist /var/lib/grafana/plugins/grafana-flightsql-datasource
