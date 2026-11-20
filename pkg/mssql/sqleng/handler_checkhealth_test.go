package sqleng

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net"
	"testing"

	"github.com/grafana/grafana-azure-sdk-go/v2/azcredentials"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/status"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrToHealthCheckResult(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want *backend.CheckHealthResult
	}{
		{
			name: "without error",
			want: &backend.CheckHealthResult{Status: backend.HealthStatusError, Message: "Internal Server Error"},
		},
		{
			name: "network error",
			err:  errors.Join(errors.New("foo"), &net.OpError{Op: "read", Net: "tcp", Err: errors.New("some op")}),
			want: &backend.CheckHealthResult{
				Status:      backend.HealthStatusError,
				Message:     "Network error: Failed to connect to the server. Error message: some op",
				JSONDetails: []byte(`{"errorDetailsLink":"https://grafana.com/docs/grafana/latest/datasources/mssql","verboseMessage":"foo\nread tcp: some op"}`),
			},
		},
		{
			name: "db error",
			err:  errors.Join(errors.New("foo"), &mssql.Error{Message: "error foo occurred in mssql server"}),
			want: &backend.CheckHealthResult{
				Status:      backend.HealthStatusError,
				Message:     "foo\nmssql: error foo occurred in mssql server (0)",
				JSONDetails: []byte(`{"errorDetailsLink":"https://grafana.com/docs/grafana/latest/datasources/mssql","verboseMessage":"foo\nmssql: error foo occurred in mssql server (0)"}`),
			},
		},
		{
			name: "regular error",
			err:  errors.New("internal server error"),
			want: &backend.CheckHealthResult{
				Status:      backend.HealthStatusError,
				Message:     "internal server error",
				JSONDetails: []byte(`{"errorDetailsLink":"https://grafana.com/docs/grafana/latest/datasources/mssql","verboseMessage":"internal server error"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ErrToHealthCheckResult(tt.err)
			require.Nil(t, err)
			assert.Equal(t, string(tt.want.JSONDetails), string(got.JSONDetails))
			require.Equal(t, tt.want, got)
		})
	}
}

type failingConnector struct{ err error }

func (c failingConnector) Connect(context.Context) (driver.Conn, error) { return nil, c.err }
func (c failingConnector) Driver() driver.Driver                        { return nil }

func TestCheckHealthErrorSource(t *testing.T) {
	adminReq := &backend.CheckHealthRequest{PluginContext: backend.PluginContext{User: &backend.User{Role: "Admin"}}}
	viewerReq := &backend.CheckHealthRequest{PluginContext: backend.PluginContext{User: &backend.User{Role: "Viewer"}}}
	noUserReq := &backend.CheckHealthRequest{}
	connErr := errors.New("socks connect tcp 10.0.0.1:1433->db:1433: unknown error host unreachable")

	tests := []struct {
		name             string
		db               *sql.DB
		req              *backend.CheckHealthRequest
		azureCurrentUser bool
		want             status.Source
	}{
		{name: "ping failure, admin", db: sql.OpenDB(failingConnector{connErr}), req: adminReq, want: status.SourceDownstream},
		{name: "ping failure, non-admin", db: sql.OpenDB(failingConnector{connErr}), req: viewerReq, want: status.SourceDownstream},
		{name: "ping failure, no user", db: sql.OpenDB(failingConnector{connErr}), req: noUserReq, want: status.SourceDownstream},
		{name: "login failure", db: sql.OpenDB(failingConnector{&mssql.Error{Number: 18456, Message: "Login failed for user"}}), req: viewerReq, want: status.SourceDownstream},
		{name: "missing Azure current user is a configuration error", req: adminReq, azureCurrentUser: true, want: status.SourceDownstream},
		{name: "connection not initialized is a plugin error", db: nil, req: adminReq, want: status.SourcePlugin},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &DataSourceHandler{db: tt.db, log: backend.NewLoggerWith("logger", "mssql.test"), queryResultTransformer: &testQueryResultTransformer{}}
			if tt.azureCurrentUser {
				e.dsInfo.JsonData.AuthenticationType = azureAuthentication
				e.azureCredentials = &azcredentials.AadCurrentUserCredentials{}
			}
			ctx := status.InitSource(context.Background())

			res, err := e.CheckHealth(ctx, tt.req)
			require.NoError(t, err)
			require.Equal(t, backend.HealthStatusError, res.Status)
			require.Equal(t, tt.want, status.SourceFromContext(ctx))
		})
	}
}
