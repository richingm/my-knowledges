#!/bin/sh

/usr/local/bin/app -conf /data/conf &

nginx -g "daemon off;"