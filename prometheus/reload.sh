#!/bin/bash
# reload the config
hostname="$1"

case $hostname in
local)
  host="localhost"
;;
*)
  host="$hostname"
;;
esac

curl -X POST "http://"$host":9090/-/reload"