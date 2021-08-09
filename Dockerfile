FROM scratch
LABEL description="Couchbase image to cleanup rotated audit logs" vendor="Couchbase" maintainer="docker@couchbase.com"

COPY non-root.passwd /etc/passwd
COPY LICENSE /licenses/couchbase.txt
COPY README.md /help.1

COPY build/bin/couchbase-audit-cleanup /usr/local/bin/couchbase-audit-cleanup
RUN chmod a+x /usr/local/bin/couchbase-audit-cleanup

USER 8453
CMD [ "/usr/local/bin/couchbase-audit-cleanup" ]