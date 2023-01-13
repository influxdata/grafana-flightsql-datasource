package flightsql

import (
	"context"
	"crypto/x509"
	"fmt"

	"github.com/apache/arrow/go/v10/arrow/flight/flightsql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// newFlightSQLClient creats a new flightsql.Client from a Datasource
// configuration structure. It establishes appropriate transport and RPC
// security settings for the gRPC connection.
func newFlightSQLClient(cfg config) (*flightsql.Client, error) {
	dialOptions, err := grpcDialOptions(cfg)
	if err != nil {
		return nil, fmt.Errorf("grpc dial options: %s", err)
	}
	return flightsql.NewClient(cfg.Host, nil, nil, dialOptions...)
}

// bearerToken is a credentials.PerRPCCredentials implementation that populates
// the `authorization` metadata field for all RPCs produced by a gRPC client.
type bearerToken struct {
	token                    string
	requireTransportSecurity bool
}

func (t bearerToken) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (t bearerToken) RequireTransportSecurity() bool {
	return t.requireTransportSecurity
}

// grpcDialOptions produces appropriate gRPC dialing options w.r.t. security for
// a given Datasource configuration structure.
func grpcDialOptions(cfg config) ([]grpc.DialOption, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
	}

	if cfg.Secure {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("x509: %s", err)
		}
		return append(opts,
			grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(pool, "")),
			grpc.WithPerRPCCredentials(bearerToken{token: cfg.Token, requireTransportSecurity: true}),
		), nil
	}
	return append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(bearerToken{token: cfg.Token}),
	), nil
}
