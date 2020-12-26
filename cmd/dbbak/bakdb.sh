#!/usr/bin/env bash
datetime=`date +%Y%m%d_%H%M%S |cut -b1-20`
mongodump -h 139.186.84.15 --port 27989 -u poolweb -p xjrw2020 -d poolweb -o ~/dbbak/$datetime/

