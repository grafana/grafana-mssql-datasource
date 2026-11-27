package schema_test

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/grafana/dsconfig/schema"
	"github.com/grafana/grafana-mssql-datasource/pkg/mssql/sqleng"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/pluginschema"
	"k8s.io/kube-openapi/pkg/spec3"
)

//go:embed dsconfig.json
var configSchemaJSON []byte

// sqleng.JsonData includes fields shared with other SQL plugins. Keep the
// fields MSSQL uses, then include the settings read outside this struct.
func jsonDataModel() any {
	fields := reflect.VisibleFields(reflect.TypeFor[struct {
		MaxIdleConnsAuto          bool    `json:"maxIdleConnsAuto"`
		KeytabFilePath            string  `json:"keytabFilePath"`
		CredentialCache           string  `json:"credentialCache"`
		CredentialCacheLookupFile string  `json:"credentialCacheLookupFile"`
		ConfigFilePath            string  `json:"configFilePath"`
		UDPConnectionLimit        int     `json:"UDPConnectionLimit"`
		EnableDNSLookupKDC        string  `json:"enableDNSLookupKDC"`
		OAuthPassThru             bool    `json:"oauthPassThru"`
		DisableGrafanaCache       bool    `json:"disableGrafanaCache"`
		Timeout                   float64 `json:"timeout"`
		KeepAlive                 float64 `json:"keepAlive"`
		AzureCredentials          struct {
			AuthType   string `json:"authType"`
			AzureCloud string `json:"azureCloud"`
			TenantID   string `json:"tenantId"`
			ClientID   string `json:"clientId"`
			UserID     string `json:"userId"`
		} `json:"azureCredentials"`
	}]())
	for _, field := range reflect.VisibleFields(reflect.TypeFor[sqleng.JsonData]()) {
		switch field.Tag.Get("json") {
		case "timescaledb", "sslmode", "tlsConfigurationMethod", "sslCertFile", "sslKeyFile", "timezone", "allowCleartextPasswords":
			continue
		case "servername":
			// The editor uses serverName. Go decodes this key case-insensitively.
			field.Tag = `json:"serverName"`
		}
		fields = append(fields, field)
	}
	return reflect.New(reflect.StructOf(fields)).Interface()
}

//go:generate go test -run TestPlugin -generateArtifacts
func TestPlugin(t *testing.T) {
	schema.RunPluginTests(t, schema.PluginUnderTest{
		ID:                "mssql",
		ConfigSchemaJSON:  configSchemaJSON,
		SettingsJSONModel: jsonDataModel(),
		SecureKeys:        []string{"password", "azureClientSecret", "clientSecret", "secureSocksProxyPassword"},
		SettingsExamples: &pluginschema.SettingsExamples{Examples: map[string]*spec3.Example{
			"": example("SQL Server Authentication", "grafana", map[string]any{
				"encrypt": "false", "maxOpenConns": 100, "maxIdleConns": 100, "maxIdleConnsAuto": true,
				"connMaxLifetime": 14400, "connectionTimeout": 0,
			}, "password"),
			"windows": exampleAt("Windows Authentication when Grafana runs on Windows under an account with SQL Server access",
				"localhost:1433", "Windows Authentication", "", nil, ""),
			"kerberosPassword": example("Windows AD: Username + password", "name@EXAMPLE.COM", nil, "password"),
			"kerberosKeytab": example("Windows AD: Keytab", "name@EXAMPLE.COM", map[string]any{
				"keytabFilePath": "/home/grot/grot.keytab",
			}, ""),
			"kerberosCache": example("Windows AD: Credential cache", "", map[string]any{
				"credentialCache": "/tmp/krb5cc_1000",
			}, ""),
			"kerberosCacheLookup": example("Windows AD: Credential cache file", "name@EXAMPLE.COM", map[string]any{
				"credentialCacheLookupFile": "/home/grot/cache.json",
			}, ""),
			"appRegistration": azureExample("App Registration", map[string]string{
				"authType": "clientsecret", "azureCloud": "AzureCloud",
				"tenantId": "00000000-0000-0000-0000-000000000001",
				"clientId": "00000000-0000-0000-0000-000000000002",
			}, "azureClientSecret"),
			"managedIdentity": azureExample("Managed Identity", map[string]string{"authType": "msi"}, ""),
			"entraPassword": azureExample("Azure Entra Password", map[string]string{
				"authType": "ad-password", "userId": "user@example.com",
				"clientId": "00000000-0000-0000-0000-000000000002",
			}, "password"),
			"currentUser": azureExample("Current User", map[string]string{"authType": "currentuser"}, ""),
			"skipTLSVerification": exampleAt("SQL Server Authentication with encryption and skipped certificate verification",
				"localhost:1433", "SQL Server Authentication", "grafana", map[string]any{
					"encrypt": "true", "tlsSkipVerify": true,
				}, "password"),
			"privateConnection": exampleAt("SQL Server Authentication via a Grafana Cloud private data source connection",
				"sql.example.com:1433", "SQL Server Authentication", "grafana", map[string]any{
					"enableSecureSocksProxy": true, "secureSocksProxyUsername": "<PRIVATE_NETWORK_ID>", "encrypt": "true",
				}, "password"),
			"strictEncryption": exampleAt("Self-managed Grafana with strict encryption and a local CA for SQL Server 2022 or later",
				"sql.example.com:1433", "SQL Server Authentication", "grafana", map[string]any{
					"encrypt": "strict", "serverName": "sql.example.com",
					"sslRootCertFile": "/etc/grafana/certs/sql-ca.pem",
				}, "password"),
		}},
	})
}

func azureExample(summary string, credentials map[string]string, secret string) *spec3.Example {
	settings := map[string]any{"azureCredentials": credentials, "encrypt": "true"}
	if credentials["authType"] == "currentuser" {
		settings["oauthPassThru"] = true
		settings["disableGrafanaCache"] = true
	}
	return exampleAt(summary, "example.database.windows.net:1433", "Azure AD Authentication", "", settings, secret)
}

func example(auth, user string, settings map[string]any, secret string) *spec3.Example {
	return exampleAt(auth, "localhost:1433", auth, user, settings, secret)
}

func exampleAt(summary, host, auth, user string, settings map[string]any, secret string) *spec3.Example {
	jsonData := map[string]any{"database": "grafana", "authenticationType": auth}
	for key, value := range settings {
		jsonData[key] = value
	}
	spec := map[string]any{"title": "MSSQL", "url": host, "jsonData": jsonData}
	if user != "" {
		spec["user"] = user
	}
	resource := map[string]any{
		"apiVersion": "mssql.datasource.grafana.app/v0alpha1",
		"kind":       "DataSource",
		"metadata":   map[string]string{"name": "mssql"},
		"spec":       spec,
	}
	if secret != "" {
		resource["secure"] = map[string]any{secret: map[string]string{"create": "<REPLACE_WITH_SECRET>"}}
	}
	return &spec3.Example{ExampleProps: spec3.ExampleProps{Summary: summary, Value: resource}}
}
