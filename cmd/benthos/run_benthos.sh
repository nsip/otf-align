#!/bin/bash


rm -r ./msgs
mkdir ./msgs

BENTHOS_HOST=localhost:34196

benthos lint workflow.yaml

benthos -c workflow.yaml