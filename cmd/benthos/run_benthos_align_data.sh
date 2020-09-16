#!/bin/bash


benthos lint alignData.yaml

clear && benthos -c alignData.yaml
