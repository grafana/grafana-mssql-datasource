//go:build integration

package kerberos_test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
)

func TestKerberosEnvVarFallback(t *testing.T) {
	t.Setenv("KRB5_CONFIG", "/krb5/crossrealm.conf")
	t.Setenv("KRB5CCNAME", "/runtime/bob-tgt.ccache")
	ds := newDS(t, "", "", map[string]any{"authenticationType": "Windows AD: Credential cache"})
	require.Equal(t, backend.HealthStatusOk, checkHealth(t, ds).Status)
	login, _ := queryLogin(t, ds)
	require.Equal(t, `A\bob`, login)
}
