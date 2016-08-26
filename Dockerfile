FROM scratch
ADD rootdev /rootdev
ADD build /build
ADD cacert.pem /etc/ssl/certs/ca-certificates.crt

EXPOSE 8022
CMD ["/rootdev"]
