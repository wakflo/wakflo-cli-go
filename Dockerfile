FROM alpine:3.20
COPY wakflo-cli /usr/bin/wakflo-cli
ENTRYPOINT ["/usr/bin/wakflo-cli"]