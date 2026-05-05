package mssql

import (
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/microsoft/go-mssqldb/integratedauth/krb5"

	"github.com/grafana/grafana-azure-sdk-go/v2/azsettings"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/config"

	"github.com/grafana/grafana-mssql-datasource/pkg/mssql/sqleng"
)

func NewInstanceSettings(logger log.Logger) datasource.InstanceFactoryFunc {
	return func(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
		grafCfg := config.GrafanaConfigFromContext(ctx)
		sqlCfg, err := grafCfg.SQL()
		if err != nil {
			return nil, err
		}
		pluginCfg := backend.PluginConfigFromContext(ctx)

		jsonData := sqleng.JsonData{
			MaxOpenConns:      sqlCfg.DefaultMaxOpenConns,
			MaxIdleConns:      sqlCfg.DefaultMaxIdleConns,
			ConnMaxLifetime:   sqlCfg.DefaultMaxConnLifetimeSeconds,
			Encrypt:           "false",
			ConnectionTimeout: 0,
			SecureDSProxy:     false,
		}

		azureSettings, err := azsettings.ReadSettings(ctx)
		if err != nil {
			logger.Error("failed to read Azure settings from Grafana", "error", err.Error())
			return nil, err
		}

		err = json.Unmarshal(settings.JSONData, &jsonData)
		if err != nil {
			return nil, fmt.Errorf("error reading settings: %w", err)
		}

		database := jsonData.Database
		if database == "" {
			database = settings.Database
		}

		dsInfo := sqleng.DataSourceInfo{
			JsonData:                jsonData,
			URL:                     settings.URL,
			User:                    settings.User,
			Database:                database,
			ID:                      settings.ID,
			Updated:                 settings.Updated,
			UID:                     settings.UID,
			DecryptedSecureJSONData: settings.DecryptedSecureJSONData,
			OrgID:                   pluginCfg.OrgID, // nolint:staticcheck
		}

		userFacingDefaultError, err := grafCfg.UserFacingDefaultError()
		if err != nil {
			return nil, err
		}

		config := sqleng.DataPluginConfiguration{
			DSInfo:            dsInfo,
			MetricColumnTypes: []string{"VARCHAR", "CHAR", "NVARCHAR", "NCHAR"},
			RowLimit:          sqlCfg.RowLimit,
		}
		handler, err := sqleng.NewQueryDataHandler(ctx, settings, userFacingDefaultError, config, logger, azureSettings)
		if err != nil {
			logger.Error("Failed connecting to MSSQL", "err", err)
			return nil, err
		}

		logger.Debug("Successfully connected to MSSQL")
		return handler, nil
	}
}
