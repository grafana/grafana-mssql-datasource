#!/bin/bash
set -euo pipefail

samba-tool user create sqlsvc Lab-Svc-Pass-1
samba-tool user setexpiry sqlsvc --noexpiry
samba-tool spn add MSSQLSvc/sql.b.test:1433 sqlsvc
samba-tool spn add MSSQLSvc/sql.b.test sqlsvc
# Must run before exportkeytab, or the SPN keys are RC4 only.
ldbmodify -H /var/lib/samba/private/sam.ldb <<'EOF'
dn: CN=sqlsvc,CN=Users,DC=b,DC=test
changetype: modify
replace: msDS-SupportedEncryptionTypes
msDS-SupportedEncryptionTypes: 24
EOF

samba-tool user create alice Lab-Alice-1
samba-tool user setexpiry alice --noexpiry

rm -f /runtime/mssql.keytab /runtime/alice.keytab
for principal in MSSQLSvc/sql.b.test:1433 MSSQLSvc/sql.b.test sqlsvc; do
  samba-tool domain exportkeytab /runtime/mssql.keytab --principal="$principal"
done
samba-tool domain exportkeytab /runtime/alice.keytab --principal=alice
chmod 444 /runtime/mssql.keytab /runtime/alice.keytab
