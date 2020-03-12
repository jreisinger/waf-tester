#!/bin/sh

GEOIP_FOLDER=/etc/nginx/geoip

geoip2_get()
{
    curl $2 -o $GEOIP_FOLDER/$1.tar.gz || { echo "Could not download $1, exiting." ; exit 1; }
    mkdir $GEOIP_FOLDER/$1 \
        && tar xf $GEOIP_FOLDER/$1.tar.gz -C $GEOIP_FOLDER/$1 --strip-components 1 \
        && mv $GEOIP_FOLDER/$1/$1.mmdb $GEOIP_FOLDER/$1.mmdb \
        && rm -rf $GEOIP_FOLDER/$1 \
        && rm -rf $GEOIP_FOLDER/$1.tar.gz
}

mkdir -p $GEOIP_FOLDER

geoip2_get "GeoLite2-City" "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz"
geoip2_get "GeoLite2-ASN" "http://geolite.maxmind.com/download/geoip/database/GeoLite2-ASN.tar.gz"
