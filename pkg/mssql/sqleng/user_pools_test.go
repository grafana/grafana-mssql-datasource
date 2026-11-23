package sqleng

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/grafana/grafana-azure-sdk-go/v2/azcredentials"
	"github.com/grafana/grafana-azure-sdk-go/v2/azsettings"
	"github.com/grafana/grafana-azure-sdk-go/v2/azusercontext"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCurrentUserHandler() *DataSourceHandler {
	return &DataSourceHandler{
		log: backend.NewLoggerWith("logger", "test"),
		dsInfo: DataSourceInfo{
			URL:      "localhost:1433",
			UID:      "ds-uid",
			OrgID:    1,
			JsonData: JsonData{AuthenticationType: azureAuthentication},
		},
		azureCredentials: &azcredentials.AadCurrentUserCredentials{},
		azureSettings: &azsettings.AzureSettings{
			UserIdentityTokenEndpoint: &azsettings.TokenEndpointSettings{ClientId: "client-id", ClientSecret: "client-secret"},
		},
		driverName:  "azuresql",
		proxyClient: proxy.New(nil),
	}
}

func withIDToken(idToken string) context.Context {
	return azusercontext.WithCurrentUser(context.Background(), azusercontext.CurrentUserContext{
		User:    &backend.User{Email: "user@example.com"},
		IdToken: idToken,
	})
}

func isClosed(db *sql.DB) bool {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(errors.New("probe the pool without dialing"))
	return !errors.Is(db.PingContext(ctx), context.Canceled)
}

func TestUserPoolIdleEviction(t *testing.T) {
	t.Run("closes and removes a pool idle past the TTL", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			h := newCurrentUserHandler()
			defer h.Dispose()
			ctx := withIDToken("token-1")
			db, err := h.getDB(ctx)
			require.NoError(t, err)

			time.Sleep(2 * userPoolIdleTTL)
			synctest.Wait()

			assert.True(t, isClosed(db))
			next, err := h.getDB(ctx)
			require.NoError(t, err)
			assert.NotSame(t, db, next)
			assert.False(t, isClosed(next))
		})
	})

	t.Run("keeps a pool used within the TTL", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			h := newCurrentUserHandler()
			defer h.Dispose()
			ctx := withIDToken("token-1")
			db, err := h.getDB(ctx)
			require.NoError(t, err)

			time.Sleep(userPoolIdleTTL - time.Minute)
			again, err := h.getDB(ctx)
			require.NoError(t, err)
			require.Same(t, db, again)

			time.Sleep(userPoolIdleTTL - time.Minute)
			synctest.Wait()

			assert.False(t, isClosed(db))
			again, err = h.getDB(ctx)
			require.NoError(t, err)
			assert.Same(t, db, again)
		})
	})
}

func TestGetDBConcurrentFirstCallsShareOnePool(t *testing.T) {
	// synctest fails the test if a pool opened by a losing caller is never closed.
	synctest.Test(t, func(t *testing.T) {
		h := newCurrentUserHandler()
		defer h.Dispose()
		ctx := withIDToken("token-1")

		const callers = 32
		dbs := make([]*sql.DB, callers)
		start := make(chan struct{})
		var wg sync.WaitGroup
		for i := range callers {
			wg.Go(func() {
				<-start
				db, err := h.getDB(ctx)
				assert.NoError(t, err)
				dbs[i] = db
			})
		}
		close(start)
		wg.Wait()

		for _, db := range dbs {
			assert.Same(t, dbs[0], db)
		}
	})
}

func TestDisposeClosesUserPools(t *testing.T) {
	t.Run("closes open pools and stops the sweeper", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			h := newCurrentUserHandler()
			first, err := h.getDB(withIDToken("token-1"))
			require.NoError(t, err)
			second, err := h.getDB(withIDToken("token-2"))
			require.NoError(t, err)

			h.Dispose()

			assert.True(t, isClosed(first))
			assert.True(t, isClosed(second))
			_, err = h.getDB(withIDToken("token-3"))
			assert.ErrorIs(t, err, errDBClosed)
		})
	})

	t.Run("before any pool opens", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			h := newCurrentUserHandler()

			h.Dispose()

			_, err := h.getDB(withIDToken("token-1"))
			assert.ErrorIs(t, err, errDBClosed)
		})
	})
}
