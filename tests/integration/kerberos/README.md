# Kerberos integration lab

This lab is a container environment for the Kerberos authentication modes of the MSSQL data source.
It runs two Samba Active Directory forests with a two-way forest trust, and CoreDNS for the lab DNS names.
Each domain controller creates its users, service principal names (SPNs), and keytabs when it provisions its domain.

## Lab values

All values are synthetic and exist only inside the lab.

| Item | Value |
|---|---|
| Network | `10.99.0.0/24`: dns `.5`, dc-b `.10`, dc-a `.20`, sql `.30`, test `.40`, mitkdc `.50`, grafana `.60`, tools `.70` |
| SQL Server realm | `B.TEST`, NetBIOS `B`, DC `dc1.b.test`, SQL Server `sql.b.test:1433` |
| User realm | `A.TEST`, NetBIOS `A`, DC `dc1.a.test`, two-way forest trust with `B.TEST` |
| Administrator (both DCs) | `Lab-Admin-1` |
| `sqlsvc@B.TEST` (privileged account, SPNs `MSSQLSvc/sql.b.test:1433`, `MSSQLSvc/sql.b.test`) | `Lab-Svc-Pass-1`, `msDS-SupportedEncryptionTypes: 24`, no expiry |
| `alice@B.TEST`, `bob@A.TEST` | `Lab-Alice-1`, `Lab-Bob-Pass-1`, no expiry |
| SQL Server `sa` | `Lab-Sa-Pass-1` |
| MIT realms | `ALPHA.TEST` to `HUB.TEST` to `GAMMA.TEST`, `bob@ALPHA.TEST` `Lab-Bob-Mit-1`, SPN `MSSQLSvc/svc.gamma.test:1433@GAMMA.TEST` |
| Runtime files | named volume `runtime` at `/runtime`: `mssql.keytab`, `alice.keytab`, `bob.keytab`, `alice-tgt.ccache`, `bob-tgt.ccache`, `dc-a.ready`, `dc-b.ready` |

## Start, stop, and reset the lab

Run these commands from the repository root.

To start the lab:

```shell
docker compose -f tests/integration/kerberos/compose.yaml up -d --build --wait
```

The command returns when both domain controllers are healthy.
`dc-b` is healthy only after the forest trust exists.

To stop the lab and keep its domains:

```shell
docker compose -f tests/integration/kerberos/compose.yaml stop
```

To reset the lab:

```shell
docker compose -f tests/integration/kerberos/compose.yaml down -v
```

`down -v` is the reset.
It deletes the containers and the `runtime` volume with its keytabs and credential caches.
The next start provisions new domains with new keys and writes new runtime files.
A keytab from an earlier lab has the wrong keys for a new lab.
After you change a setup script, or when a start fails, reset the lab.

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
