FROM scratch
ADD rootdev /rootdev
ADD build /build

EXPOSE 8022
CMD ["/rootdev"]
