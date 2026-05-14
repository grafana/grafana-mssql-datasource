import { test, expect } from '@grafana/plugin-e2e';

import type { MssqlOptions } from '../../src/types';

const PLUGIN_TYPE = 'mssql';
const PROVISIONED_FILE = 'datasources.yml';

function exploreUrl(uid: string, opts?: { rawSql?: string; format?: string }) {
  const query: Record<string, unknown> = {
    refId: 'A',
    datasource: { type: PLUGIN_TYPE, uid },
    format: opts?.format ?? 'table'
  };
  if (opts?.rawSql) {
    query.rawSql = opts.rawSql;
  }
  const panes = JSON.stringify({
    explore: {
      datasource: uid,
      queries: [query],
      range: { from: 'now-4h', to: 'now' },
    },
  });
  return `/explore?orgId=1&schemaVersion=1&panes=${encodeURIComponent(panes)}`;
}

async function switchEditorMode(page: Page, mode: 'Builder' | 'Code') {
  const target = page.getByRole('radio', { name: mode, exact: true });
  if (await target.isChecked()) {
    return;
  }
  await target.click();
  if (mode === 'Builder') {
    const discardButton = page.getByRole('button', { name: 'Discard code and switch' })
    if (await discardButton.isVisible({ timeout: 3000 }).catch(() => false)) {
      await discardButton.click();
    }
  }
  await expect(target).toBeChecked();
}

test.describe('Query editor', () => {
  test.describe('rendering', () => {
    test(
      'smoke: renders all editor mode options',
      { tag: '@plugins' },
      async ({ page, readProvisionedDataSource }) => {
        const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
        await page.goto(exploreUrl(ds.uid));
        await expect(page.getByRole('radio', { name: 'Builder', exact: true })).toBeVisible();
        await expect(page.getByRole('radio', { name: 'Code', exact: true })).toBeVisible();
      }
    );

    test('renders format dropdown in all modes', async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await page.goto(exploreUrl(ds.uid));
      const formatSelect = page.getByTestId('query-editor-row').getByRole('combobox', { name: 'Format' });
      await switchEditorMode(page, 'Builder');
      await expect(formatSelect).toBeVisible();
      await switchEditorMode(page, 'Code');
      await expect(formatSelect).toBeVisible();
    });

    test('renders run query button in all modes', async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await page.goto(exploreUrl(ds.uid));
      const runQueryButton = page.getByTestId('query-editor-row').getByRole('button', { name: 'Run query', exact: true });
      await switchEditorMode(page, 'Builder');
      await expect(runQueryButton).toBeVisible();
      await switchEditorMode(page, 'Code');
      await expect(runQueryButton).toBeVisible();
    });
  });

  test.describe('builder mode', () => {
    test('shows expected fields', async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await page.goto(exploreUrl(ds.uid));
      await switchEditorMode(page, 'Builder');

      await expect(page.getByRole('switch', { name: 'Filter' })).toBeVisible();
      await expect(page.getByRole('switch', { name: 'Group' })).toBeVisible();
      await expect(page.getByRole('switch', { name: 'Order' })).toBeVisible();
      await expect(page.getByRole('switch', { name: 'Preview' })).toBeVisible();

      await expect(page.getByRole('combobox', { name: 'Dataset' })).toBeVisible();
      await expect(page.getByRole('combobox', { name: 'Table' })).toBeVisible();
      await expect(page.getByRole('combobox', { name: 'Data operations' })).toBeVisible();
      await expect(page.getByRole('combobox', { name: 'Column' })).toBeVisible();
      await expect(page.getByRole('combobox', { name: 'Alias' })).toBeVisible();
    });

    test('can select a dataset, table, and column', async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await page.goto(exploreUrl(ds.uid));
      await switchEditorMode(page, 'Builder');

      const queryRow = page.getByTestId('query-editor-row');

      await page.getByRole('combobox', { name: 'Dataset' }).click();
      await page.getByRole('option', { name: 'grafanadb' }).click();
      await expect(queryRow).toContainText('grafanadb');

      await page.getByRole('combobox', { name: 'Table' }).click();
      await page.getByRole('option', { name: 'dbo.WORLD_DATA' }).click();
      await expect(queryRow).toContainText('dbo.WORLD_DATA');

      await page.getByRole('combobox', { name: 'Column' }).click();
      await page.getByRole('option', { name: 'BASE_COUNTRY' }).click();
      await expect(queryRow).toContainText('BASE_COUNTRY');
    });
  });

  test.describe('code mode', () => {
    test('shows expected fields', async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await page.goto(exploreUrl(ds.uid));
      await switchEditorMode(page, 'Code');
      await expect(page.getByRole('textbox', { name: /editor content/i })).toBeVisible();
    });
  
    test('can enter a SQL query string', async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await page.goto(exploreUrl(ds.uid));
      await switchEditorMode(page, 'Code');
      const textbox = page.getByRole('textbox', { name: /editor content/i });
      await textbox.fill('SELECT TOP 5 * FROM dbo.WORLD_DATA');
      await expect(textbox).toHaveValue('SELECT TOP 5 * FROM dbo.WORLD_DATA');
    });
  });
});

test.describe('Query editor with fixture data', () => {
  test.describe.configure({ mode: 'serial' });

  test.describe('Afghanistan', () => {
    test('table query returns rows', async ({ page, explorePage, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      const sql = "SELECT * FROM dbo.WORLD_DATA WHERE JSON_VALUE(BASE_COUNTRY, '$.Name') = 'Afghanistan'";
      let body: Record<string, unknown> | null = null;
      const responsePromise = explorePage.waitForQueryDataResponse(async (r) => {
        if (!r.ok()) {
          return false;
        }
        const b: any = await r.json().catch(() => null);
        if (!Array.isArray(b?.results?.A?.frames)) {
          return false;
        }
        body = b;
        return true;
      });
      await page.goto(exploreUrl(ds.uid, { rawSql: sql }));
      await responsePromise;
      expect((body as any)?.results?.A?.frames?.length).toBeGreaterThan(0);
    });

    test('query returns expected column names', async ({ page, explorePage, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      const sql = "SELECT TOP 1 * FROM dbo.WORLD_DATA WHERE JSON_VALUE(BASE_COUNTRY, '$.Name') = 'Afghanistan'";
      let body: Record<string, unknown> | null = null;
      const responsePromise = explorePage.waitForQueryDataResponse(async (r) => {
        if (!r.ok()) {
          return false;
        }
        const b: any = await r.json().catch(() => null);
        if (!Array.isArray(b?.results?.A?.frames)) {
          return false;
        }
        body = b;
        return true;
      });
      await page.goto(exploreUrl(ds.uid, { rawSql: sql }));
      await responsePromise;
      const fields: string[] = (body as any)?.results?.A?.frames?.[0]?.schema?.fields?.map(
        (f: { name: string }) => f.name
      );
      expect(fields).toContain('BASE_COUNTRY');
      expect(fields).toContain('BIRTH_RATE');
      expect(fields).toContain('CO2');
      expect(fields).toContain('GDP');
      expect(fields).toContain('DATE_TIME');
      expect(fields).toContain('TIMESTAMP_VALUE');
    });
  });

  test.describe('Albania', () => {
    test('table query returns rows', async ({ page, explorePage, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      const sql = "SELECT * FROM dbo.WORLD_DATA WHERE JSON_VALUE(BASE_COUNTRY, '$.Name') = 'Albania'";
      let body: Record<string, unknown> | null = null;
      const responsePromise = explorePage.waitForQueryDataResponse(async (r) => {
        if (!r.ok()) {
          return false;
        }
        const b: any = await r.json().catch(() => null);
        if (!Array.isArray(b?.results?.A?.frames)) {
          return false;
        }
        body = b;
        return true;
      });
      await page.goto(exploreUrl(ds.uid, { rawSql: sql }));
      await responsePromise;
      expect((body as any)?.results?.A?.frames?.length).toBeGreaterThan(0);
    });

    test('query with time range filter returns rows', async ({ page, explorePage, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      const sql = "SELECT DATE_TIME as time, GDP FROM dbo.WORLD_DATA WHERE $__timeFilter(DATE_TIME) AND JSON_VALUE(BASE_COUNTRY, '$.Name') = 'Albania'";
      let body: Record<string, unknown> | null = null;
      const responsePromise = explorePage.waitForQueryDataResponse(async (r) => {
        if (!r.ok()) {
          return false;
        }
        const b: any = await r.json().catch(() => null);
        if (!Array.isArray(b?.results?.A?.frames)) {
          return false;
        }
        body = b;
        return true;
      });
      await page.goto(exploreUrl(ds.uid, { rawSql: sql }));
      await responsePromise;
      expect((body as any)?.results?.A?.frames?.length).toBeGreaterThan(0);
    });
  });
});
