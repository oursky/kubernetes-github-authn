import json
from fabric.operations import run, local
from fabric.api import env, settings
from fabric.decorators import serial,runs_once

env.shell = "/bin/sh -l -c"

# Test locally
# docker build --no-cache -t k8sauthn .
# docker run -it --rm -p 8080:3000 k8sauthn

def build():
    local('GOOS=linux GOARCH=amd64 go build -o _output/main main.go')
    local('docker build --no-cache -t k8sauthn .')
    local('rm -rf _output')
