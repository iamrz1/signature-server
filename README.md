# Signature Server

Signature server stores transaction blobs and uses predefined secret key to sign and verify those transactions.

## How to run
Signature server doesn't have any external dependency and can run on most machines out of the box.
This server has been tested on MacOS 12.1 Monterey, using go version go1.17.2 darwin/amd64.

Build and run this server:
```shell
$ make build
$ make run
```
This will create a binary and run it on local machine.

The server can be also run without creating an explicit binary:
```shell
$ make test_run
```

In any case, a config.yml is used as the configuration file. It contains the port numbers, and the secret key.
For kubernetes deployment, a sample deployment yaml is available in the deployment directory.

Test cases has also been added for the APIs that can be run with the following command:
```shell
$ make test 
```

## API Doc:

This Signature server includes swagger for api documentation. It is exposed at `/doc/` endpoint (ie. http://localhost:8080/doc/index.html#/ for a local instance).
