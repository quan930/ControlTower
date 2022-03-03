FROM docker:20.10.12-dind

COPY server_dind.sh /
COPY dockerd-entrypoint2.sh /usr/local/bin/
RUN apk add --no-cache socat && apk add --no-cache bash
VOLUME /var/lib/docker
EXPOSE 2375 2376

ENTRYPOINT ["dockerd-entrypoint2.sh"]