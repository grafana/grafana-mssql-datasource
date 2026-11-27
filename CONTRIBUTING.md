# Contributing

## Signed commits are required

> [!IMPORTANT]
> All commits must be [signed](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits) (GPG, SSH, or S/MIME) to be merged into this repository. Pull requests with unsigned commits will need to be re-committed with signatures before they can be merged.

Thank you for your interest in contributing to the Microsoft SQL Server data source for Grafana! We welcome contributions from the community.

Feel free to [browse open issues](https://github.com/grafana/grafana-mssql-datasource/issues) or open a new one. For more general guidance, see [Grafana's Contributing Guide](https://github.com/grafana/grafana/blob/main/CONTRIBUTING.md).

This project adheres to the [Grafana Code of Conduct](https://github.com/grafana/grafana/blob/main/CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Prerequisites

- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/) (see [go.mod](go.mod) for the minimum required version)
- [Mage](https://magefile.org/)
- [Node.js LTS](https://nodejs.org)
- [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm) (see [package.json](package.json) for the minimum required version)
- [Docker](https://docs.docker.com/get-docker/)

## Frontend

1. Install dependencies:

   ```shell
   npm install
   ```

2. Build plugin in development mode and watch for changes:

   ```shell
   npm run dev
   ```

3. Build plugin in production mode:

   ```shell
   npm run build
   ```

4. Run frontend tests:

   ```shell
   npm run test:ci
   ```

## Backend

1. Build the backend binaries:

   ```shell
   mage -v
   ```

## Local development environment

`npm run server` starts a local Microsoft SQL Server instance and a Grafana instance with the plugin pre-provisioned:

```shell
npm run server
```

## Kerberos integration tests

The Kerberos tests run against a container lab with two Active Directory forests, SQL Server, and an MIT KDC:

```shell
docker compose -f tests/integration/kerberos/compose.yaml --profile test run --build --rm test go test -tags integration ./tests/integration/...
```

To browse the lab in Grafana, build the plugin (`npm run build` and `mage -v`) and run:

```shell
npm run server:kerberos
```

Grafana is on `http://localhost:3001`. See [the lab README](tests/integration/kerberos/README.md) for realms, principals, and passwords, and how to stop and reset it.

Keep these points in mind when you change the lab:

- SQL Server needs reverse DNS that returns the domain controller FQDN, so keep the CoreDNS PTR records and do not use Docker's DNS.
- SQL Server on Linux needs the NetBIOS name of each domain (`B`, `A`) to resolve to its domain controller, with the FQDN first. Keep these entries in `sql/hosts`, which is bind-mounted on `/etc/hosts`, because Docker reorders `extra_hosts`.
- SQL Server looks up Windows logins with the privileged AD account `sqlsvc` (`network.privilegedadaccount` in `sql/mssql.conf`). Its password must never expire.
- The SQL Server keytab must be owned by `mssql:mssql` with mode 440 and hold both SPN forms, `MSSQLSvc/host:port` and `MSSQLSvc/host`. Set the AES encryption types on `sqlsvc` before you export the keytab, or the SPN keys are RC4 only.
- A client that asks for a service in another realm needs an exact `host:port` entry in `[domain_realm]`, for example `sql.b.test:1433 = B.TEST` in `krb5/crossrealm.conf`. The lab's Samba and MIT KDCs send no referrals. The go-krb5 library looks up the service realm with the port attached, so the `.b.test` suffix rule does not match. Windows AD sends referrals, so it needs no mapping.
- Credential caches must hold a TGT only. A cached service ticket skips the cross-realm exchange that the tests need to cover.
- Never run `sqlservr` or `mssql-conf` as root in the `sql` container.
- Grafana loads the external plugin only when `[plugin.mssql] as_external = true` is in an ini file. The environment variable does not work.
- `[plugins] forward_host_env_vars = mssql` passes `KRB5_CONFIG` and `KRB5CCNAME` to the plugin process.

## E2E tests

```shell
npm run server
npm run e2e
```

Or, to install Playwright browsers first:

```shell
npx playwright install --with-deps
npm run server
npm run e2e
```

## Release

You need commit access to the repository to publish a release.

1. Update the version number in `package.json`.
2. Update `CHANGELOG.md` with the changes included in the release.
3. Open a PR with the changes and merge it.
4. Follow the release process described [here](https://enghub.grafana-ops.net/docs/default/component/grafana-plugins-platform/plugins-ci-github-actions/010-plugins-ci-github-actions/#cd_1).
