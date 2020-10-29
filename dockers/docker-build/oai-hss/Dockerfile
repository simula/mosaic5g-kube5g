# Create an image that host oai-cn
FROM ubuntu:16.04
LABEL maintainer="Osama Arouk (arouk@eurecom.fr)"
ENV container docker
ENV PATH "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin"
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y fuse snapd snap-confine squashfuse sudo &&\
  apt-get autoremove -y &&\
 apt-get autoclean -y &&\
 dpkg-divert --local --rename --add /sbin/udevadm &&\
 ln -s /bin/true /sbin/udevadm

# Install helper tools
RUN apt-get install -y apt-utils dnsutils net-tools iputils-ping vim tcpdump
RUN mv /usr/sbin/tcpdump /usr/bin/

# Use rc.local to help
RUN sed -i "13s/^/\/root\/init.sh/g" /etc/rc.local
RUN chmod +x /etc/rc.local
COPY ./init_build.sh /root/init.sh
COPY ./init_deploy.sh /root/
COPY ./hook /root/
RUN mkdir /root/config
COPY ./conf.yaml /root/config/

# Start snapd
RUN systemctl enable snapd

# stop agetty that uses a lot of resources
RUN rm -f /lib/systemd/system/systemd*udev* \
  && rm -f /lib/systemd/system/getty.target

VOLUME ["/sys/fs/cgroup"]

CMD ["/sbin/init"]
