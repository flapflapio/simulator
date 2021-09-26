# asset-service [![Build](https://github.com/atomata/model-service/actions/workflows/test.yml/badge.svg)](https://github.com/atomata/model-service/actions/workflows/test.yml)

## QuickStart

The segment below contains commands for common tasks that you can run. You don't
need to execute all the commands, you can pick and choose depending on what you
want to do.

```bash
# Building the app:
make build

# Running the built executable:
./app

# Building the docker container:
make docker

# Running the docker container in development mode:
make docker-run
```

## How to add new APIs

This service uses controllers to group APIs together. To add more APIs, simply
write a new controller. There are 2 steps to add a new controller:

1. Create the controller in `core/controllers` directory.
2. Attach the controller to the application by creating a controller instance
   and placing it the `controllers` slice in `app.go`. The controllers in this
   slice will get added to the app:
   ```go
   // Add any new controllers to this slice
   controllers = []types.Controller{
     assetcontroller.
       New(storageService).
       WithPrefix("/assets"),
   }
   ```
3. (Optional) There is also a `core/services` directory where you can put
   services that you create, or you can put the services next to the controller
   files - either way is fine.

**NOTE:** controllers should implement the `core/types/Controller` interface. It
only has 1 function: `Attach(router Router)` which is called by the server
object to add the controllers routes to the app instance.

## Docker

**Dockerfile:** [`docker/Dockerfile`](docker/Dockerfile)

The Docker image uses a multi-stage build - a build stage is used to create the
Go binary and then the final image is built by copying the executable onto the
[`scratch`](https://hub.docker.com/_/scratch) docker image. This is a special
docker image that contains basically nothing.

Because our base image contains nothing, we are linking the app statically - so
there is no dependency on glibc.
