#!/bin/bash


benthos lint alignData.yaml

clear && benthos --chilled -c alignData.yaml
