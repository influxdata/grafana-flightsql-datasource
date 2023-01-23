package flightsql

import (
	"crypto/x509"
	"fmt"

	"github.com/apache/arrow/go/v10/arrow/flight/flightsql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func newFlightSQLClient(cfg config) (*flightsql.Client, error) {
	dialOptions, err := grpcDialOptions(cfg)
	if err != nil {
		return nil, fmt.Errorf("grpc dial options: %s", err)
	}
	return flightsql.NewClient(cfg.Addr, nil, nil, dialOptions...)
}

func grpcDialOptions(cfg config) ([]grpc.DialOption, error) {
	transport := grpc.WithTransportCredentials(insecure.NewCredentials())
	if cfg.Secure {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("x509: %s", err)
		}
		transport = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(pool, ""))
	}

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		transport,
	}

	return opts, nil
}
