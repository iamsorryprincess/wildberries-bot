#!/bin/sh

systemctl stop wbdpapi.service
cp api /usr/local/bin/wbdpapi
rm api
systemctl start wbdpapi.service
systemctl status wbdpapi.service
