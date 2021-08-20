#!/bin/bash

# ensure the necessary context has been created on the
# n3 server
curl -s  -X POST http://localhost:1323/admin/newdemocontext -d userName=nsipOtf -d contextName=alignmentMaps


# now run the workflow
# 
benthos lint alignMaps.yaml
clear && benthos --chilled -c alignMaps.yaml
