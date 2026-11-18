#!/bin/bash
set -euo pipefail
for r in ALPHA HUB GAMMA; do
  kdb5_util create -s -r $r.TEST -P Lab-Master-1 >/dev/null
done
k() { kadmin.local -r "$1" -q "$2" >/dev/null; }
k ALPHA.TEST "addprinc -pw Lab-Bob-Mit-1 bob@ALPHA.TEST"
k ALPHA.TEST "addprinc -pw Lab-Trust-1 krbtgt/HUB.TEST@ALPHA.TEST"
k HUB.TEST   "addprinc -pw Lab-Trust-1 krbtgt/HUB.TEST@ALPHA.TEST"
k HUB.TEST   "addprinc -pw Lab-Trust-2 krbtgt/GAMMA.TEST@HUB.TEST"
k GAMMA.TEST "addprinc -pw Lab-Trust-2 krbtgt/GAMMA.TEST@HUB.TEST"
k GAMMA.TEST "addprinc -randkey MSSQLSvc/svc.gamma.test:1433@GAMMA.TEST"
