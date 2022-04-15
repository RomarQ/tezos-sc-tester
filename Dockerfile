FROM alpine:3.14

# Create application directory
RUN mkdir -p /opt/tezos

ADD ./tezos-bin/tezos-client    /opt/tezos/tezos-client
ADD ./bin/linux_amd64/api       /opt/tezos/tezos-testing
ADD ./docker/config             /opt/tezos/config

EXPOSE 5000/tcp

ENTRYPOINT ["/opt/tezos/tezos-testing", "-config=/opt/tezos/config/api.yaml"]
