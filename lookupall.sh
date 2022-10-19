#!/bin/sh
for i in {2..50}; do curl h"ttp://172.17.0.$i:8080/objects/$1"; done
