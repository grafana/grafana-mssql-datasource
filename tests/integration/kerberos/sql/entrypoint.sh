#!/bin/bash
set -euo pipefail

until [ -f /runtime/mssql.keytab ]; do sleep 1; done
mkdir -p /var/opt/mssql/secrets
install -m 440 /runtime/mssql.keytab /var/opt/mssql/secrets/mssql.keytab
exec /opt/mssql/bin/launch_sqlservr.sh /opt/mssql/bin/sqlservr
