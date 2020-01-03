# nmuxj

## Story

That is a simple go rest api that aim to be a kind of webhook proxy for
the cloud-init web calls from the VM's when they are ready for the first 
time provisioning.

Then all the vm's calling this simple api and when it reach all the expected calls
then its perform a jenkins remote call to trigger a defined job with the provided
job parameters.

This make it possible to run provision jobs with jenkins as part of a cloud init
run.

## Introduction

The name nmuxj stand for nth multiplexer jenkins. That's what it is or what it does it receive a defined number of input requests and sends a build request to jenkins.

## Structure

* jenkins
    * contains jenkins docker container with job as code that configure two example jobs (provision-redis, provision-elasticsearch)
* service
    * contains golang source code file main.go
* test
    * contains the json payload to test a jenkins remote job execution with and without build parameter's 

## API

POST /provision/{provisionID}/group/{group}/total/{total:[0-9]+}/vm/{vmID}
```json
{
    "ctx":"provisioning",
    "job":"provision-elasticsearch",  // job and github repo name in jenkins
    "jobParams":                // all the build barameter needed by the deploy job
    {
        "revision":"origin/master",  //branch to use for the jenkins job
        "SERVER":"elasticsearch-1",
        "PROVISION_TARGET":"install"
    }
}
```

## Build

```bash
make build
```

## Run

## TODO

[x] build jenkins docker container with bootstrap job
[x] create the example jenkins jobs with bootstrap job
