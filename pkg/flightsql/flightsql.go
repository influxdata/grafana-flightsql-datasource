package flightsql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/apache/arrow/go/v10/arrow/flight/flightsql"
	"github.com/go-chi/chi/v5"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"google.golang.org/grpc/metadata"
)

var (
	_ backend.QueryDataHandler      = (*FlightSQLDatasource)(nil)
	_ backend.CheckHealthHandler    = (*FlightSQLDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*FlightSQLDatasource)(nil)
	_ backend.CallResourceHandler   = (*FlightSQLDatasource)(nil)
)

const mdBucketName = "bucket-name"

type config struct {
	Addr     string            `json:"host"`
	Database string            `json:"database"`
	Metadata map[string]string `json:"metadata"`
	Secure   bool              `json:"secure"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Token    string            `json:"token"`
}

func (cfg config) validate() error {
	if strings.Count(cfg.Addr, ":") == 0 {
		return fmt.Errorf(`server address must be in the form "host:port"`)
	}
	if len(cfg.Database) == 0 {
		return fmt.Errorf("database name is required")
	}

	noToken := len(cfg.Token) == 0
	noUserPass := len(cfg.Username) == 0 || len(cfg.Password) == 0

	if noToken && noUserPass {
		return fmt.Errorf("token or username/password are required")
	}

	return nil
}

// FlightSQLDatasource is a Grafana datasource plugin for Flight SQL.
type FlightSQLDatasource struct {
	database        string
	client          *flightsql.Client
	resourceHandler backend.CallResourceHandler
	md              metadata.MD
}

// NewDatasource creates a new datasource instance.
func NewDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var cfg config
	err := json.Unmarshal(settings.JSONData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("config: %s", err)
	}
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation: %v", err)
	}

	client, err := newFlightSQLClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("flightsql: %s", err)
	}

	md := metadata.MD{}
	for k, v := range cfg.Metadata {
		md.Set(k, v)
	}
	// TODO(brett): Remove when the UI is completely configurable via metadata.
	md.Set(mdBucketName, cfg.Database)

	ctx := context.Background()
	if len(cfg.Username) > 0 || len(cfg.Password) > 0 {
		ctx, err = client.Client.AuthenticateBasicToken(ctx, cfg.Username, cfg.Password)
		if err != nil {
			return nil, fmt.Errorf("flightsql: %s", err)
		}
		authMD, _ := metadata.FromOutgoingContext(ctx)
		md = metadata.Join(md, authMD)
	} else {
		md.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Token))
	}

	ds := &FlightSQLDatasource{
		database: cfg.Database,
		client:   client,
		md:       md,
	}
	r := chi.NewRouter()
	r.Use(recoverer)
	r.Route("/plugin", func(r chi.Router) {
		r.Get("/macros", ds.getMacros)
	})
	r.Route("/flightsql", func(r chi.Router) {
		r.Get("/sql-info", ds.getSQLInfo)
		r.Get("/tables", ds.getTables)
		r.Get("/columns", ds.getColumns)
	})
	ds.resourceHandler = httpadapter.New(r)

	return ds, nil
}

// Dispose cleans up before we are reaped.
func (d *FlightSQLDatasource) Dispose() {
	if err := d.client.Close(); err != nil {
		logErrorf(err.Error())
	}
}

// CallResource forwards requests to an internal HTTP mux that handles custom
// resources for the datasource.
func (d *FlightSQLDatasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return d.resourceHandler.CallResource(ctx, req, sender)
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *FlightSQLDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if resp := d.query(ctx, "select 1"); resp.Error != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("ERROR: %s", resp.Error),
		}, nil
	}
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "OK",
	}, nil
}

func recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if rec == http.ErrAbortHandler {
					panic(rec)
				}
				logErrorf("Panic: %v %v", r, string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func logInfof(format string, v ...any) {
	log.DefaultLogger.Info(fmt.Sprintf(format, v...))
}

func logErrorf(format string, v ...any) {
	log.DefaultLogger.Error(fmt.Sprintf(format, v...))
}
