FROM scratch
ADD rootdev /rootdev
ADD build /build
ADD cacert.pem /etc/ssl/certs/ca-certificates.crt

LABEL traefik.domain="rootdev.nl"
EXPOSE 8022
CMD ["/rootdev"]
