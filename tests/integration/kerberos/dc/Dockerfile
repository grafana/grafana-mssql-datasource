FROM debian:trixie-slim
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
      samba samba-ad-dc samba-ad-provision samba-dsdb-modules samba-vfs-modules \
      winbind krb5-user ldb-tools smbclient dnsutils iproute2 procps ca-certificates \
 && rm -rf /var/lib/apt/lists/*
COPY entrypoint.sh setup-a.sh setup-b.sh /
RUN chmod +x /entrypoint.sh /setup-a.sh /setup-b.sh
ENTRYPOINT ["/entrypoint.sh"]
