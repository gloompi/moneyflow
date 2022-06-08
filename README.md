# Moneyflow API

[![CircleCI](https://circleci.com/gh/ardanlabs/service.svg?style=svg)](https://circleci.com/gh/ardanlabs/service) [![Go Report Card](https://goreportcard.com/badge/github.com/ardanlabs/service)](https://goreportcard.com/report/github.com/ardanlabs/service) [![go.mod Go version](https://img.shields.io/github/go-mod/go-version/ardanlabs/service)](https://github.com/ardanlabs/service)

Copyright 2022, Kubanychbek Esenzhanov gloompi@gmail.com


## Commands

`make service` - builds a docker container with executable binary inside

`make kind-up` - creates a cluster and configures a namespace for a project
    
`make kind-down` - deletes a cluster created by kind-up
    
`make kind-load` - loads docker image into the cluster created by kind-up
    
`make kind-apply` - builds a kustomized configuration, creates a pod and applies all the changes to a pod that being created
    
`make kind-status` - will show all the related statuses
    
`make kind-status-service` - will show active pods in our namespace
    
`make kind-logs` - will show golang logs
    
`make kind-restart` - restarts our pod
    
`make kind-update` - updates our app after changes made to our source code
    
`make kind-update-apply` - updates a pod after changing configurations
    
`make kind-describe` - will show detailed information about pod, including error logs
    
`make tidy` - will install dependencies and create a vendor folder

## Techonologies

[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/ardanlabs/service)](https://github.com/ardanlabs/service)

Kubernetes

Docker

Kind

Kustomize