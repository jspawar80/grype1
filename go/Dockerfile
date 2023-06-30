FROM alpine:latest

RUN apk update && apk add --no-cache \
    curl \
    docker-cli \
    jq \
    bash

RUN LATEST_GRYPE=$(curl --silent "https://api.github.com/repos/anchore/grype/releases/latest" | jq -r '.tag_name') && \
    curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh | sh -s -- -b /usr/local/bin ${LATEST_GRYPE}

RUN apk del curl jq

VOLUME /var/run/docker.sock:/var/run/docker.sock

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD [ "node", "server.js" ]

