# Create an image that host cassandra
#  docker build -t mosaic5gecosys/cassandra-m5g:1.1 --force-rm=true --rm=true .  |& tee build.log

FROM mosaic5gecosys/cassandra-m5g:1.0
LABEL maintainer="Osama Arouk (arouk@eurecom.fr)"
ENV container docker


RUN cp /usr/local/bin/docker-entrypoint.sh /usr/local/bin/docker-entrypoint_tmp.sh
RUN sed '3a /root/hook &' /usr/local/bin/docker-entrypoint_tmp.sh  > /usr/local/bin/docker-entrypoint_tmp2.sh
RUN cp /usr/local/bin/docker-entrypoint_tmp2.sh /usr/local/bin/docker-entrypoint.sh

COPY ./hook /root/
RUN touch /root/hook.log
