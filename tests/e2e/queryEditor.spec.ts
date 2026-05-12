import { test, expect } from '@grafana/plugin-e2e';

import type { MssqlOptions } from '../../src/types';

const PLUGIN_TYPE = 'mssql';
const PROVISIONED_FILE = 'datasources.yml';

// Fixture data constants — seed=42
// 24 rows per country (Afghanistan, Albania) at 10-minute intervals
const FIXTURE_FROM_ISO = '2026-05-11T20:00:00.000Z';
const FIXTURE_TO_ISO = '2026-05-12T00:00:00.000Z';

function exploreUrl(uid: string, opts?: { rawSql?: string; format?: string }) {
  const query: Record<string, unknown> = {
    refId: 'A',
    datasource: { type: PLUGIN_TYPE, uid },
    format: opts?.format ?? 'table',
    editorMode: 'code',
  };
  if (opts?.rawSql) {
    query.rawSql = opts.rawSql;
  }
  const panes = JSON.stringify({
    explore: {
      datasource: uid,
      queries: [query],
      range: { from: FIXTURE_FROM_ISO, to: FIXTURE_TO_ISO },
    },
  });
  return `/explore?orgId=1&schemaVersion=1&panes=${encodeURIComponent(panes)}`;
}

test.describe('Query editor', () => {
  test.describe('rendering', () => {
    test('smoke: renders Builder and Code editor mode options', { tag: '@plugins' }, async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await page.goto(exploreUrl(ds.uid));
      await expect(page.getByRole('radio', { name: 'Builder' })).toBeVisible();
      await expect(page.getByRole('radio', { name: 'Code' })).toBeVisible();
    });

    test('renders Format dropdown in Code mode', async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await page.goto(exploreUrl(ds.uid));
      await expect(page.getByRole('combobox', { name: /Format:/ })).toBeVisible();
    });
  });

  test.describe('Code mode', () => {
    test('shows Monaco editor with pre-populated SQL', async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      const sql = 'SELECT TOP 5 * FROM dbo.WORLD_DATA';
      await page.goto(exploreUrl(ds.uid, { rawSql: sql }));
      await expect(page.getByRole('radio', { name: 'Code' })).toBeChecked();
      await expect(page.getByRole('textbox', { name: /editor content/i })).toHaveValue(sql);
    });
  });

  test.describe('Builder mode', () => {
    test('shows Dataset and Table dropdowns after switching from Code', async ({ page, readProvisionedDataSource }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await page.goto(exploreUrl(ds.uid));
      await page.getByRole('radio', { name: 'Builder' }).click();
      const discardButton = page.getByRole('button', { name: 'Discard code and switch' });
      if (await discardButton.isVisible({ timeout: 2000 }).catch(() => false)) {
        await discardButton.click();
      }
      await expect(page.getByRole('radio', { name: 'Builder' })).toBeChecked();
      await expect(page.getByRole('combobox', { name: /Dataset/ })).toBeVisible();
      await expect(page.getByRole('combobox', { name: /Table/ })).toBeVisible();
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
      const sql = `SELECT DATE_TIME as time, GDP FROM dbo.WORLD_DATA WHERE DATE_TIME >= '2026-05-11 20:00:00' AND DATE_TIME < '2026-05-12 00:00:00' AND JSON_VALUE(BASE_COUNTRY, '$.Name') = 'Albania'`;
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
