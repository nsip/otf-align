#!/bin/bash


rm -r ./msgs
mkdir ./msgs

benthos lint align_workflow.yaml

clear && benthos -c align_workflow.yaml