FROM alpine:3.21
COPY wakflo-cli /usr/bin/wakflo-cli
ENTRYPOINT ["/usr/bin/wakflo-cli"]