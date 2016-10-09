FROM scratch
MAINTAINER mpdroog <rootdev@gmail.com>
ADD rootdev /rootdev
ADD build /build
ADD cacert.pem /etc/ssl/certs/ca-certificates.crt

ARG GIT_COMMIT=unknown
LABEL git-commit=$GIT_COMMIT 

LABEL traefik.backend=rootdev
LABEL traefik.frontend.rule=Host:rootdev.nl

EXPOSE 8080
ENTRYPOINT ["/rootdev"]
