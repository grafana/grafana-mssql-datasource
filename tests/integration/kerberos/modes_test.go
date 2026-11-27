//go:build integration

package kerberos_test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
)

func TestKerberosModes(t *testing.T) {
	cases := []struct {
		name, user, password, wantLogin string
		jsonData                        map[string]any
	}{
		{"b-password", "alice@B.TEST", "Lab-Alice-1", `B\alice`, map[string]any{"authenticationType": "Windows AD: Username + password", "configFilePath": "/krb5/single.conf"}},
		{"b-keytab", "alice@B.TEST", "", `B\alice`, map[string]any{"authenticationType": "Windows AD: Keytab", "keytabFilePath": "/runtime/alice.keytab", "configFilePath": "/krb5/single.conf"}},
		{"b-ccache", "", "", `B\alice`, map[string]any{"authenticationType": "Windows AD: Credential cache", "credentialCache": "/runtime/alice-tgt.ccache", "configFilePath": "/krb5/single.conf"}},
		{"b-lookup", "alice", "", `B\alice`, map[string]any{"authenticationType": "Windows AD: Credential cache file", "credentialCacheLookupFile": "/fixtures/lookup.json", "configFilePath": "/krb5/single.conf"}},
		{"a-password", "bob@A.TEST", "Lab-Bob-Pass-1", `A\bob`, map[string]any{"authenticationType": "Windows AD: Username + password", "configFilePath": "/krb5/crossrealm.conf"}},
		{"a-keytab", "bob@A.TEST", "", `A\bob`, map[string]any{"authenticationType": "Windows AD: Keytab", "keytabFilePath": "/runtime/bob.keytab", "configFilePath": "/krb5/crossrealm.conf"}},
		{"a-ccache", "", "", `A\bob`, map[string]any{"authenticationType": "Windows AD: Credential cache", "credentialCache": "/runtime/bob-tgt.ccache", "configFilePath": "/krb5/crossrealm.conf"}},
		{"a-lookup", "bob", "", `A\bob`, map[string]any{"authenticationType": "Windows AD: Credential cache file", "credentialCacheLookupFile": "/fixtures/lookup.json", "configFilePath": "/krb5/crossrealm.conf"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ds := newDS(t, tc.user, tc.password, tc.jsonData)
			health := checkHealth(t, ds)
			require.Equal(t, backend.HealthStatusOk, health.Status, health.Message)
			login, scheme := queryLogin(t, ds)
			require.Equal(t, tc.wantLogin, login)
			require.Equal(t, "KERBEROS", scheme)
		})
	}
}

func TestKerberosWrongPassword(t *testing.T) {
	ds := newDS(t, "alice@B.TEST", "not-the-password", map[string]any{"authenticationType": "Windows AD: Username + password", "configFilePath": "/krb5/single.conf"})
	health := checkHealth(t, ds)
	require.Equal(t, backend.HealthStatusError, health.Status)
	require.Contains(t, health.Message, "KDC_ERR_PREAUTH_FAILED")
}
