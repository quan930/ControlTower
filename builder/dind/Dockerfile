FROM docker:20.10.12-dind

COPY ter.sh /
COPY dockerd-entrypoint2.sh /usr/local/bin/

VOLUME /var/lib/docker
EXPOSE 2375 2376

ENTRYPOINT ["dockerd-entrypoint2.sh"]