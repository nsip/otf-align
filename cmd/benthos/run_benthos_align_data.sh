#!/bin/bash


rm -r ./msgs
mkdir ./msgs

benthos lint alignData.yaml

clear && benthos -c alignData.yaml
