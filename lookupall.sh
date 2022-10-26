#!/bin/bash
for i in {2..51}; do curl "http://172.17.0.$i:8080/objects/$1"; done
