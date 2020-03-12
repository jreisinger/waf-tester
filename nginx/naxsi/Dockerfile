# Based on https://github.com/daxio/nginx-naxsi/blob/master/Dockerfile

# Nginx version here should probably match NGINX_VERSION below
FROM nginx:1.16.0

ENV DEBIAN_FRONTEND=noninteractive
ENV NGINX_VERSION=1.16.0
ENV NAXSI_VERSION=0.55.3

# Install basic packages and build tools
RUN apt-get update && \
    apt-get install --no-install-recommends --no-install-suggests -y \
      curl \
      wget \
      build-essential \
      libssl-dev \
      libpcre3 \
      libpcre3-dev \
      zlib1g-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Download and extract sources
RUN cd /tmp && \
    wget --no-check-certificate \
        https://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz \
        -O nginx.tar.gz && \
    wget --no-check-certificate \
        https://github.com/nbs-system/naxsi/archive/${NAXSI_VERSION}.tar.gz \
        -O naxsi.tar.gz && \
    tar -xvf nginx.tar.gz && \
    mv nginx-${NGINX_VERSION} nginx && \
    rm nginx.tar.gz && \
    tar -xvf naxsi.tar.gz && \
    mv naxsi-${NAXSI_VERSION} naxsi && \
    rm naxsi.tar.gz && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* nginx.tar.gz.asc naxsi.tar.gz.asc

# Configure and build
RUN cd /tmp/nginx && \
    BASE_CONFIGURE_ARGS=`nginx -V 2>&1 | grep "configure arguments" | cut -d " " -f 3-` && \
    /bin/sh -c "./configure --add-module=/tmp/naxsi/naxsi_src ${BASE_CONFIGURE_ARGS}" && \
    make && make install && \
    rm -rf /tmp/nginx*

# Forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log

COPY nginx.conf /etc/nginx/
COPY naxsi_core.rules /etc/nginx/
COPY naxsi_custom.rules /etc/nginx/
COPY naxsi.conf /etc/nginx/
COPY naxsi_block_page.conf /etc/nginx/

#HEALTHCHECK --interval=10s --timeout=3s CMD http_proxy="" curl --fail http://localhost/healthcheck || exit 1

#EXPOSE 1234

STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]
