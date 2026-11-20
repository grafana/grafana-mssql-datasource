#!/bin/bash
set -euo pipefail

samba-tool user create bob Lab-Bob-Pass-1
samba-tool user setexpiry bob --noexpiry

rm -f /runtime/bob.keytab
samba-tool domain exportkeytab /runtime/bob.keytab --principal=bob
chmod 444 /runtime/bob.keytab
