I'm using this https://github.com/bxcodec/go-clean-arch repository as a reference

### How To Run This Project

Since the project already use Go Module, I recommend to put the source code in any folder but GOPATH.

#### Run the Applications
Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/bxcodec/go-clean-arch.git

#move to project
$ cd go-clean-arch

# Build the docker image first
$ make docker

# Run the application
$ make run

# check if the containers are running
$ docker ps

# Execute the call
$ curl localhost:9090/products

# Stop
$ make stop
```
