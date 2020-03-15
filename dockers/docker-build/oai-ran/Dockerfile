# Create an image that host oai-ran
FROM ubuntu:16.04
ENV container docker
ENV PATH "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin"
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y fuse snapd snap-confine squashfuse sudo &&\
 apt-get clean &&\
 dpkg-divert --local --rename --add /sbin/udevadm &&\
 ln -s /bin/true /sbin/udevadm

RUN apt-get install -y apt-utils
RUN apt-get install -y dnsutils
RUN apt-get install -y net-tools
RUN apt-get install -y iputils-ping
RUN apt-get install -y vim

# Use rc.local to help
RUN sed -i "13s/^/\/root\/init.sh/g" /etc/rc.local
RUN chmod +x /etc/rc.local
COPY ./init.sh /root/
COPY ./hook /root/
RUN mkdir /root/config
COPY ./conf.yaml /root/config/

# Start snapd
RUN systemctl enable snapd

VOLUME ["/sys/fs/cgroup"]
STOPSIGNAL SIGRTMIN+3
CMD ["/sbin/init"]