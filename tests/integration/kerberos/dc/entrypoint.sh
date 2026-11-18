#!/bin/bash
set -euo pipefail
: "${REALM:?}" "${DOMAIN:?}" "${ADMIN_PASSWORD:?}" "${TGT_USER:?}"

dc="dc-${DOMAIN,,}"
rm -f "/runtime/$dc.ready"

fresh=false
if [ ! -f /var/lib/samba/private/sam.ldb ]; then
  fresh=true
  rm -f /etc/samba/smb.conf
  samba-tool domain provision \
    --realm="$REALM" --domain="$DOMAIN" --server-role=dc \
    --dns-backend=SAMBA_INTERNAL --adminpass="$ADMIN_PASSWORD" \
    --option="dns forwarder = 10.99.0.5"
  "/setup-${DOMAIN,,}.sh"
fi

samba --interactive --no-process-group --debug-stdout &
samba_pid=$!

ccache="/runtime/$TGT_USER-tgt.ccache"
get_tgt() { kinit -k -t "/runtime/$TGT_USER.keytab" -c "FILE:$ccache" "$TGT_USER@$REALM"; }
tries=120
until get_tgt 2>/dev/null; do
  kill -0 "$samba_pid"
  tries=$((tries - 1))
  if [ "$tries" -le 1 ]; then
    get_tgt || exit 1
    break
  fi
  sleep 1
done
chmod 444 "$ccache"

if [ "$fresh" = true ] && [ "$dc" = dc-b ]; then
  until [ -f /runtime/dc-a.ready ]; do sleep 2; done
  samba-tool domain trust create a.test --type=forest --direction=both --create-location=both \
    -U Administrator@A.TEST --password="$ADMIN_PASSWORD" \
    --local-dc-username=Administrator --local-dc-password="$ADMIN_PASSWORD"
fi

touch "/runtime/$dc.ready"
wait "$samba_pid"
