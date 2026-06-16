import { test, expect } from '@grafana/plugin-e2e';

import type { MssqlOptions } from '../../src/types';

const PLUGIN_TYPE = 'mssql';
const PROVISIONED_FILE = 'datasources.yml';

const DS_URL = process.env.DS_INSTANCE_URL ?? 'mssql:1433';
const DS_DATABASE = process.env.DS_INSTANCE_DATABASE ?? 'grafanadb';
const DS_AUTHN_TYPE = process.env.DS_INSTANCE_AUTHN_TYPE ?? 'SQL Server Authentication';
const DS_USERNAME = process.env.DS_INSTANCE_USERNAME ?? 'sa';

test.describe('Config editor', () => {
  test.describe('rendering', () => {
    test('smoke: should render config editor', { tag: '@plugins' }, async ({
      createDataSourceConfigPage,
      page
    }) => {
      await createDataSourceConfigPage({ type: PLUGIN_TYPE });
      await expect(page.getByText(/^Type\s*Microsoft SQL Server$/, { exact: true })).toBeVisible();
      await expect(page.getByRole('heading', { name: 'Connection', exact: true })).toBeVisible();
    });

    test('should render Connection section', async ({ createDataSourceConfigPage, page }) => {
      await createDataSourceConfigPage({ type: PLUGIN_TYPE });
      const section = page.getByRole('heading', { name: 'Connection', exact: true });
      await section.scrollIntoViewIfNeeded();
      await expect(page.getByPlaceholder('localhost:1433')).toBeVisible();
      await expect(page.getByPlaceholder('database name')).toBeVisible();
    });

    test('should render TLS/SSL Auth section', async ({ createDataSourceConfigPage, page }) => {
      await createDataSourceConfigPage({ type: PLUGIN_TYPE });
      const section = page.getByRole('heading', { name: 'TLS/SSL Auth', exact: true });
      await section.scrollIntoViewIfNeeded();
      await expect(section).toBeVisible();
      await expect(page.getByRole('combobox', { name: /Encrypt/, exact: true })).toBeVisible();
    });

    test('should render Authentication section', async ({ createDataSourceConfigPage, page }) => {
      await createDataSourceConfigPage({ type: PLUGIN_TYPE });
      const section = page.getByRole('heading', { name: 'Authentication', exact: true });
      await section.scrollIntoViewIfNeeded();
      await expect(section).toBeVisible();
      await expect(page.getByRole('combobox', { name: /Authentication Type/, exact: true })).toBeVisible();
      await expect(page.getByPlaceholder('user')).toBeVisible();
    });
  });

  test.describe('provisioned datasource', () => {
    test('should load provisioned connection settings', async ({
      readProvisionedDataSource,
      gotoDataSourceConfigPage,
      page,
    }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await gotoDataSourceConfigPage(ds.uid);
      await expect(page.getByPlaceholder('localhost:1433')).toHaveValue(DS_URL);
      await expect(page.getByPlaceholder('database name')).toHaveValue(DS_DATABASE);
    });

    test('should load provisioned authentication settings', async ({
      readProvisionedDataSource,
      gotoDataSourceConfigPage,
      page,
    }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await gotoDataSourceConfigPage(ds.uid);
      // The Field for Authentication Type links its label via for="authenticationType".
      // The parent div contains both the label and the Select, so toContainText finds the
      // selected option text regardless of react-select's internal DOM structure.
      await expect(page.locator('label[for="authenticationType"]').locator('xpath=parent::*')).toContainText(DS_AUTHN_TYPE);
      if (DS_AUTHN_TYPE === 'SQL Server Authentication') {
        await expect(page.getByPlaceholder('user')).toHaveValue(DS_USERNAME);
        await expect(page.getByPlaceholder('Password')).toHaveValue('configured');
      }
    });
  });

  test.describe('save & test', () => {
    test('should pass health check for provisioned datasource', async ({
      readProvisionedDataSource,
      gotoDataSourceConfigPage,
      page,
    }) => {
      const ds = await readProvisionedDataSource<MssqlOptions>({ fileName: PROVISIONED_FILE });
      await gotoDataSourceConfigPage(ds.uid);
      await page.getByRole('button', { name: /^(Save & test|Test)$/ }).click();
      await expect(page.getByText('Database Connection OK')).toBeVisible({ timeout: 15000 });
    });

    test('should show error alert when health check fails', async ({ createDataSourceConfigPage, page }) => {
      await createDataSourceConfigPage({ type: PLUGIN_TYPE });
      await page.getByPlaceholder('localhost:1433').fill(DS_URL);
      await page.route(/\/api\/datasources\/uid\/.*\/health$/, (route) =>
        route.fulfill({
          status: 400,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'connection refused' }),
        })
      );
      await page.getByRole('button', { name: /^(Save & test|Test)$/ }).click();
      await expect(page.getByText('connection refused')).toBeVisible({ timeout: 10000 });
    });

    test('should show error alert when backend is unreachable', async ({ createDataSourceConfigPage, page }) => {
      await createDataSourceConfigPage({ type: PLUGIN_TYPE });
      await page.getByPlaceholder('localhost:1433').fill(DS_URL);
      await page.route(/\/api\/datasources\/uid\/.*\/health$/, (route) =>
        route.fulfill({
          status: 503,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'dial tcp unreachable-host:1433: no such host' }),
        })
      );
      await page.getByRole('button', { name: /^(Save & test|Test)$/ }).click();
      await expect(page.getByText('no such host')).toBeVisible({ timeout: 10000 });
    });
  });
});
