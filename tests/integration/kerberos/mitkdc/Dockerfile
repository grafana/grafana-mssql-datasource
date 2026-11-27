FROM debian:trixie-slim
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends krb5-kdc krb5-admin-server krb5-user && rm -rf /var/lib/apt/lists/*
COPY kdc.conf /etc/krb5kdc/kdc.conf
COPY krb5.conf /etc/krb5.conf
COPY init.sh /init.sh
RUN chmod +x /init.sh && /init.sh
CMD ["krb5kdc", "-n", "-r", "ALPHA.TEST", "-r", "HUB.TEST", "-r", "GAMMA.TEST"]
