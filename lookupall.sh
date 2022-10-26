#!/bin/sh
for i in {2..51}; do curl h"ttp://172.17.0.$i:8080/objects/$1"; done
