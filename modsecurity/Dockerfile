FROM alpine:latest

ENV NGINX_VERSION=1.17.9 \
    NGINX_CONFIG="\
    --prefix=/usr/share/nginx \
    --sbin-path=/usr/sbin/nginx \
    --conf-path=/etc/nginx/nginx.conf \
    --modules-path=/etc/nginx/modules \
    --http-log-path=/var/log/nginx/access.log \
    --error-log-path=/var/log/nginx/error.log \
    --lock-path=/var/run/nginx.lock \
    --pid-path=/run/nginx.pid \
    --http-client-body-temp-path=/var/cache/nginx/body \
    --http-fastcgi-temp-path=/var/cache/nginx/fastcgi \
    --http-proxy-temp-path=/var/cache/nginx/proxy \
    --http-scgi-temp-path=/var/cache/nginx/scgi \
    --http-uwsgi-temp-path=/var/cache/nginx/uwsgi \
    --user=nginx \
    --group=nginx \
    --with-compat \
    --with-pcre-jit \
    --with-http_ssl_module \
    --with-http_stub_status_module \
    --with-http_realip_module \
    --with-http_auth_request_module \
    --with-http_addition_module \
    --with-http_dav_module \
    --with-http_geoip_module \
    --with-http_gzip_static_module \
    --with-http_sub_module \
    --with-http_v2_module \
    --with-stream \
    --with-stream_ssl_module \
    --with-stream_ssl_preread_module \
    --with-threads \
    --with-http_secure_link_module \
    --with-http_gunzip_module \
    --without-mail_pop3_module \
    --without-mail_smtp_module \
    --without-mail_imap_module \
    --without-http_uwsgi_module \
    --without-http_scgi_module \
    --add-module=/usr/src/ModSecurity-nginx\
    " \
    MODSECURITY_LIB_VERSION=3.0.4 \
    OWASP_MODSECURITY_CRS_VERSION=3.2.0 \
    MODSECURITY_CONFIG="\
    --disable-doxygen-doc \
    --disable-doxygen-html \
    --disable-examples \
    --disable-dependency-tracking \
    " \
    NGX_MODSECURITY_VERSION="fix/auditlog" \
    NGX_MODSECURITY_REPO="https://github.com/AirisX/ModSecurity-nginx"

RUN addgroup -S nginx \
    && adduser -D -S -h /var/cache/nginx -s /sbin/nologin -G nginx nginx \
    # Install NGINX & ModSecurity build dependencies
    && apk add --no-cache --virtual .build-deps \
        git \
        curl \
        binutils \
        gcc \
        g++ \
        make \
        autoconf \
        automake \
        bison \
        flex \
        libtool \
        linux-headers \
        geoip-dev \
        libc-dev \
        libxml2-dev \
        openssl-dev \
        pcre-dev \
        yajl-dev \
        zlib-dev \
    # Prepare directory for Build
    && mkdir -p /usr/src \
    # Download NGINX
    && curl -fSL https://nginx.org/download/nginx-$NGINX_VERSION.tar.gz -o /tmp/nginx.tar.gz \
    && tar -zxC /usr/src -f /tmp/nginx.tar.gz \
    && rm /tmp/nginx.tar.gz \
    # Build ModSecurity
    && cd /usr/src \
    && git clone --depth 1 -b v$MODSECURITY_LIB_VERSION --single-branch https://github.com/SpiderLabs/ModSecurity \
    && cd ModSecurity \
    && git submodule init \
    && git submodule update \
    && ./build.sh \
    && ./configure $MODSECURITY_CONFIG \
    && make -j$(getconf _NPROCESSORS_ONLN) \
    && make install \
    && mkdir -p /etc/nginx/modsecurity \
    && cp modsecurity.conf-recommended /etc/nginx/modsecurity/modsecurity.conf \
    && sed -i 's/SecRuleEngine DetectionOnly/SecRuleEngine On/' /etc/nginx/modsecurity/modsecurity.conf \
    && cp unicode.mapping /etc/nginx/modsecurity/unicode.mapping \
    # Download OWASP modsecurity crs
    && cd /usr/src \
    && git clone --depth 1 -b v$OWASP_MODSECURITY_CRS_VERSION --single-branch https://github.com/SpiderLabs/owasp-modsecurity-crs \
    && cd owasp-modsecurity-crs \
    && mkdir -p /etc/nginx/modsecurity/crs/ \
    && mv crs-setup.conf.example /etc/nginx/modsecurity/crs/setup.conf \
    && mv rules/REQUEST-900-EXCLUSION-RULES-BEFORE-CRS.conf.example rules/REQUEST-900-EXCLUSION-RULES-BEFORE-CRS.conf \
    && mv rules/RESPONSE-999-EXCLUSION-RULES-AFTER-CRS.conf.example rules/RESPONSE-999-EXCLUSION-RULES-AFTER-CRS.conf \
    && mv rules /etc/nginx/modsecurity/crs/ \
    && cd .. \
    && echo $'Include /etc/nginx/modsecurity/modsecurity.conf\n\
Include /etc/nginx/modsecurity/crs/setup.conf\n\
Include /etc/nginx/modsecurity/crs/rules/*.conf' > /etc/nginx/modsecurity/nginx-modsecurity.conf \
    # ModSecurity-nginx connector (bridge between nginx and libmodsecurity)
    && cd /usr/src \
    && git clone --depth 1 -b $NGX_MODSECURITY_VERSION --single-branch $NGX_MODSECURITY_REPO \
    # Build NGINX
    && cd /usr/src/nginx-$NGINX_VERSION \
    && ./configure $NGINX_CONFIG \
    && make -j$(getconf _NPROCESSORS_ONLN) \
    && make install \
    && mkdir /etc/nginx/conf.d \
    # Strip nginx and modsecurity.so
    && strip /usr/sbin/nginx \
    && strip /usr/sbin/nginx /usr/local/modsecurity/lib/libmodsecurity.so \
    # Cleaning
    && cd / \
    && apk del --purge .build-deps \
    && rm -rf /usr/src \
    && rm -rf /usr/local/modsecurity/lib/libmodsecurity.a \
    && rm -rf /etc/nginx/*.default \
    # Check NGINX + ModSecurity runtime dependencies and install them
    && runDeps="$( \
            scanelf --needed --nobanner --format '%n#p' /usr/sbin/nginx /usr/local/modsecurity/lib/*.so \
            | tr ',' '\n' \
            | sort -u \
            | awk 'system("[ -e /usr/local/modsecurity/lib/" $1 " ]") == 0 { next } { print "so:" $1 }' \
        )" \
    && apk add --no-cache --virtual .waf-rundeps curl $runDeps

COPY nginx.conf /etc/nginx/

# forward request and error logs to docker log collector
RUN ln -sf /dev/stdout /var/log/nginx/access.log && \
    ln -sf /dev/stderr /var/log/nginx/error.log && \
    ln -sf /dev/stdout /var/log/modsec_audit.log

#EXPOSE 80
STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]