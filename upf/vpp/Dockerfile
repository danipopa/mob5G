FROM ubuntu:20.04

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

# Disable service start during package installation
RUN echo '#!/bin/sh\nexit 101' > /usr/sbin/policy-rc.d && chmod +x /usr/sbin/policy-rc.d

# Install prerequisites and FD.io repository
RUN apt-get update && \
    apt-get install -y curl gnupg software-properties-common && \
    curl -fsSL https://packagecloud.io/fdio/release/gpgkey | apt-key add - && \
    echo "deb https://packagecloud.io/fdio/release/ubuntu focal main" > /etc/apt/sources.list.d/99fdio.list && \
    apt-get update && \
    apt-get install -y iproute2 net-tools pciutils vim wget python3 python3-pip libnuma-dev tar

#RUN wget http://ftp.debian.org/debian/pool/main/o/openssl/libssl1.1_1.1.1n-0+deb10u3_amd64.deb && \
#    dpkg -i libssl1.1_1.1.1n-0+deb10u3_amd64.deb && \
#    apt-get install -f

# Install DPDK 
#RUN apt-get install -y dpdk && apt-get clean


RUN apt-get install libunwind8

COPY deb/*.deb /root/

WORKDIR /app

# Install VPP
RUN dpkg --force-not-root --no-triggers --unpack \
    /root/libvppinfra_25.02-rc0~209-gd11d31653_amd64.deb \
    /root/vpp_25.02-rc0~209-gd11d31653_amd64.deb \
    /root/vpp-plugin-core_25.02-rc0~209-gd11d31653_amd64.deb \
    /root/vpp-plugin-devtools_25.02-rc0~209-gd11d31653_amd64.deb \
    /root/vpp-plugin-dpdk_25.02-rc0~209-gd11d31653_amd64.deb

RUN apt-get install -f -y

RUN rm -rf /root/*.deb

# Restore policy-rc.d to default
RUN rm /usr/sbin/policy-rc.d

# Copy VPP configuration
COPY vpp.conf /etc/vpp/vpp.conf

# Set VPP configuration path
ENV VPP_CONFIG=/etc/vpp/vpp.conf

# Start VPP
CMD ["vpp", "-c", "/etc/vpp/vpp.conf"]

