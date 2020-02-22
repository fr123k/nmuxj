# Usage of jocker

## Introduction

This repository is using the [jocker](https://github.com/fr123k/jocker) docker image
to start a jenkins master for testing the jenkins job remote execution.

* `provision-redis`

The test job that simulate the provision of two redis servers.

* `provision-elasticsearch`

The test job that simulate the provision of one elasticsearch server.

## Usage

```bash
make jocker
```
[Jenkins](http://localhost:8888/)
