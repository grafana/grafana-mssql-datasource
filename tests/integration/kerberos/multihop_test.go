//go:build integration

package kerberos_test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/microsoft/go-mssqldb/integratedauth"
	_ "github.com/microsoft/go-mssqldb/integratedauth/krb5"
	"github.com/microsoft/go-mssqldb/msdsn"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana-mssql-datasource/pkg/mssql/kerberos"
)

func TestKerberosTwoHopServiceTicket(t *testing.T) {
	requireLab(t)

	connStr := kerberos.Krb5ParseAuthCredentials("svc.gamma.test", "1433", "master", "bob@ALPHA.TEST", "Lab-Bob-Mit-1",
		kerberos.KerberosAuth{ConfigFilePath: "/krb5/multihop.conf", UDPConnectionLimit: 1}, log.DefaultLogger)
	cfg, err := msdsn.Parse(connStr + "serverspn=MSSQLSvc/svc.gamma.test:1433;")
	require.NoError(t, err)
	auth, err := integratedauth.GetIntegratedAuthenticator(cfg)
	require.NoError(t, err)
	require.NotNil(t, auth)
	defer auth.Free()
	token, err := auth.InitialBytes()
	require.NoError(t, err)
	require.NotEmpty(t, token)
}
