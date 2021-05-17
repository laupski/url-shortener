# URL-Shortener [![CircleCI](https://circleci.com/gh/laupski/url-shortener.svg?style=svg)](https://circleci.com/gh/laupski/url-shortener) [![Go Report Card](https://goreportcard.com/badge/github.com/laupski/url-shortener)](https://goreportcard.com/report/github.com/laupski/url-shortener)  [![codecov](https://codecov.io/gh/laupski/url-shortener/branch/main/graph/badge.svg?token=E3KGPV2A7M)](https://codecov.io/gh/laupski/url-shortener)

Proof of concept URL shortener project using Go, etcd and docker.

## Requirements
* `go` installed
* `docker` installed

## Initial Setup
To run out of the box, do the following:
* `cp example-config.toml config.toml`
* `docker-compose up`
* Navigate to localhost:8000 and the web page should appear.
