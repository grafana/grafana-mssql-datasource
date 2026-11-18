# Kerberos integration lab

This lab is a container environment for the Kerberos authentication modes of the MSSQL data source.
It runs two Samba Active Directory forests with a two-way forest trust, SQL Server 2022 in `B.TEST`, and CoreDNS for the lab DNS names.
Each domain controller creates its users, service principal names (SPNs), and keytabs when it provisions its domain.

## Lab values

All values are synthetic and exist only inside the lab.

| Item                                                                                         | Value                                                                                                                                                |
| -------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| Network                                                                                      | `10.99.0.0/24`: dns `.5`, dc-b `.10`, dc-a `.20`, sql `.30`, test `.40`, mitkdc `.50`, grafana `.60`, tools `.70`                                    |
| SQL Server realm                                                                             | `B.TEST`, NetBIOS `B`, DC `dc1.b.test`, SQL Server `sql.b.test:1433`                                                                                 |
| User realm                                                                                   | `A.TEST`, NetBIOS `A`, DC `dc1.a.test`, two-way forest trust with `B.TEST`                                                                           |
| Administrator (both DCs)                                                                     | `Lab-Admin-1`                                                                                                                                        |
| `sqlsvc@B.TEST` (privileged account, SPNs `MSSQLSvc/sql.b.test:1433`, `MSSQLSvc/sql.b.test`) | `Lab-Svc-Pass-1`, `msDS-SupportedEncryptionTypes: 24`, no expiry                                                                                     |
| `alice@B.TEST`, `bob@A.TEST`                                                                 | `Lab-Alice-1`, `Lab-Bob-Pass-1`, no expiry                                                                                                           |
| SQL Server `sa`                                                                              | `Lab-Sa-Pass-1`                                                                                                                                      |
| MIT realms                                                                                   | `ALPHA.TEST` to `HUB.TEST` to `GAMMA.TEST`, `bob@ALPHA.TEST` `Lab-Bob-Mit-1`, SPN `MSSQLSvc/svc.gamma.test:1433@GAMMA.TEST`                          |
| Runtime files                                                                                | named volume `runtime` at `/runtime`: `mssql.keytab`, `alice.keytab`, `bob.keytab`, `alice-tgt.ccache`, `bob-tgt.ccache`, `dc-a.ready`, `dc-b.ready` |

## Start, stop, and reset the lab

Run these commands from the repository root.

To start the lab:

```shell
docker compose -f tests/integration/kerberos/compose.yaml run --build --rm sql-init
```

The command returns when SQL Server has the Windows logins `B\alice` and `A\bob`.
SQL Server starts only after both domain controllers are healthy.
`dc-b` is healthy only after the forest trust exists.
Do not use `up --wait` to start the whole lab, because it fails when `sql-init` exits.
On an Apple Silicon Mac, SQL Server runs under Rosetta, and a cold start takes about 25 seconds when the images are already pulled and built.
The `test` profile runs the tests in the lab (see `CONTRIBUTING.md`), and `npm run server:kerberos` starts the lab with the `dev` profile (see "Run Grafana against the lab").

To stop the lab and keep its domains:

```shell
docker compose -f tests/integration/kerberos/compose.yaml --profile '*' stop
```

To reset the lab:

```shell
docker compose -f tests/integration/kerberos/compose.yaml --profile '*' down -v
```

`--profile '*'` makes both commands include the `dev` and `test` services, however you started the lab.
`down -v` is the reset.
It deletes the containers and the `runtime` volume with its keytabs and credential caches.
The next start provisions new domains with new keys and writes new runtime files.
A keytab from an earlier lab has the wrong keys for a new lab.
The TGT credential caches (`alice-tgt.ccache`, `bob-tgt.ccache`) expire after 10 hours, so `down -v` followed by a new start refreshes them.
After you change a setup script, or when a start fails, reset the lab.

## Run Grafana against the lab

Build the plugin (`npm run build` and `mage -v`), then start the lab with Grafana and the `tools` container:

```shell
npm run server:kerberos
```

Grafana runs on `http://localhost:3001`, so it can run at the same time as `npm run server` on port 3000.
It loads the plugin from `dist/` as an external plugin and provisions the data sources `b-password`, `b-keytab`, `b-ccache`, `b-lookup`, `a-password`, `a-keytab`, `a-ccache`, `a-lookup`, and `env-ccache`.
The `b-*` sources sign in as `alice@B.TEST` and the `a-*` sources as `bob@A.TEST`.
`env-ccache` has no Kerberos settings, so it uses the `KRB5_CONFIG` and `KRB5CCNAME` environment variables of the Grafana container.
Anonymous access with the Admin role is on, and it is only for this local lab.

## Inspect tickets with the tools container

The `tools` container has the Kerberos client commands and the same `/runtime`, `/krb5`, and `/fixtures` mounts, all read-only.
It writes its tickets to the default cache, not to `/runtime`:

```shell
docker compose -f tests/integration/kerberos/compose.yaml --profile dev exec tools bash
export KRB5_CONFIG=/krb5/crossrealm.conf
echo Lab-Bob-Pass-1 | kinit bob@A.TEST
klist
kvno MSSQLSvc/sql.b.test:1433@B.TEST
klist
```

The second `klist` shows the referral ticket `krbtgt/B.TEST@A.TEST` and the service ticket for SQL Server.
Use `KRB5_CONFIG=/krb5/single.conf` and `alice@B.TEST` for the single-realm path.

## MIT multi-hop realms

The `mitkdc` service runs three MIT realms in one KDC: `ALPHA.TEST`, `HUB.TEST` and `GAMMA.TEST`.
The realm path from `ALPHA.TEST` to `GAMMA.TEST` has two cross-realm hops through `HUB.TEST`.
Samba and AD forest trusts have one hop only, so only this path shows an authenticator realm error in the driver.
`TestKerberosTwoHopServiceTicket` signs in as `bob@ALPHA.TEST` with `krb5/multihop.conf`.
It asks the driver for the ticket for `MSSQLSvc/svc.gamma.test:1433@GAMMA.TEST`.
The KDC creates its databases and principals when the image builds, so no health check is needed.

## Add a principal

The setup scripts run only when a domain controller provisions its domain.
Samba is not running at that time, so use only `samba-tool` commands that work on the local database.

1. Add the `samba-tool user create` and `samba-tool domain exportkeytab` lines to `dc/setup-b.sh` for `B.TEST`, or to `dc/setup-a.sh` for `A.TEST`.
   For example:

   ```shell
   samba-tool user create carol Lab-Carol-1
   samba-tool user setexpiry carol --noexpiry
   samba-tool domain exportkeytab /runtime/carol.keytab --principal=carol
   chmod 444 /runtime/carol.keytab
   ```

1. Reset the lab and start it again.
