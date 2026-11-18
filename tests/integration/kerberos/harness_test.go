//go:build integration

package kerberos_test

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana-mssql-datasource/pkg/mssql"
)

const labAddress = "sql.b.test:1433"

type kerberosDS interface {
	backend.CheckHealthHandler
	backend.QueryDataHandler
}

func requireLab(t *testing.T) {
	t.Helper()
	conn, err := net.DialTimeout("tcp", labAddress, 5*time.Second)
	if err != nil {
		t.Fatalf("Kerberos lab not reachable at %s: %v\nRun the tests inside the lab with:\ndocker compose -f tests/integration/kerberos/compose.yaml --profile test run --rm test go test -tags integration ./tests/integration/...", labAddress, err)
	}
	_ = conn.Close()
}

func grafanaCtx() context.Context {
	cfg := backend.NewGrafanaCfg(map[string]string{
		"GF_USER_FACING_DEFAULT_ERROR":             "please inspect Grafana server log for details",
		"GF_SQL_ROW_LIMIT":                         "1000000",
		"GF_SQL_MAX_OPEN_CONNS_DEFAULT":            "100",
		"GF_SQL_MAX_IDLE_CONNS_DEFAULT":            "100",
		"GF_SQL_MAX_CONN_LIFETIME_SECONDS_DEFAULT": "14400",
	})
	return backend.WithGrafanaConfig(context.Background(), cfg)
}

func newDS(t *testing.T, user, password string, jsonData map[string]any) kerberosDS {
	t.Helper()
	requireLab(t)

	merged := map[string]any{"database": "master", "encrypt": "disable"}
	for k, v := range jsonData {
		merged[k] = v
	}
	raw, err := json.Marshal(merged)
	require.NoError(t, err)

	settings := backend.DataSourceInstanceSettings{
		URL:      labAddress,
		User:     user,
		JSONData: raw,
	}
	if password != "" {
		settings.DecryptedSecureJSONData = map[string]string{"password": password}
	}

	instance, err := mssql.NewInstanceSettings(log.DefaultLogger)(grafanaCtx(), settings)
	require.NoError(t, err)
	ds, ok := instance.(kerberosDS)
	require.True(t, ok, "instance does not implement kerberosDS")
	if d, ok := instance.(interface{ Dispose() }); ok {
		t.Cleanup(d.Dispose)
	}
	return ds
}

func checkHealth(t *testing.T, ds kerberosDS) *backend.CheckHealthResult {
	t.Helper()
	// CheckHealth reads User.Role without a nil check.
	res, err := ds.CheckHealth(grafanaCtx(), &backend.CheckHealthRequest{
		PluginContext: backend.PluginContext{User: &backend.User{Role: "Admin"}},
	})
	require.NoError(t, err)
	return res
}

func queryLogin(t *testing.T, ds kerberosDS) (login, scheme string) {
	t.Helper()
	q, err := json.Marshal(map[string]any{
		"rawSql": "SELECT SUSER_SNAME() AS login, auth_scheme FROM sys.dm_exec_connections WHERE session_id = @@SPID",
		"format": "table",
		"refId":  "A",
	})
	require.NoError(t, err)

	resp, err := ds.QueryData(grafanaCtx(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{User: &backend.User{Role: "Admin"}},
		Queries:       []backend.DataQuery{{RefID: "A", JSON: q, TimeRange: backend.TimeRange{From: time.Now().Add(-time.Hour), To: time.Now()}}},
	})
	require.NoError(t, err)
	r, ok := resp.Responses["A"]
	require.True(t, ok, "no response for refId A")
	require.NoError(t, r.Error)
	require.Len(t, r.Frames, 1)
	require.Len(t, r.Frames[0].Fields, 2)

	l, ok := r.Frames[0].Fields[0].ConcreteAt(0)
	require.True(t, ok)
	s, ok := r.Frames[0].Fields[1].ConcreteAt(0)
	require.True(t, ok)
	return l.(string), s.(string)
}
