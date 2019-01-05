# RUN

Get dependency first, need https://github.com/golang/dep

```bash
make deps
```
### Local Run

To run services locally
```bash
make run
```

curl api
```bash
# Company A
curl -v -H "token:35216c9e-dea4-458c-babd-325f2ef0eefb" localhost:8081/data

# Company B
curl -v -H "token:35216c9e-dea4-458c-babd-325f2ef0eefb" localhost:8082/data

# Public Car Rental API
curl -v localhost:8083/data
```
### Docker Run

run build. Make sure your stop local run to free up ports
```bash
# build binary and docker images
make build
``` 

Start containers
```bash
# Sart containers
make run-containers
```

expected services
```bash
$ docker ps -a
CONTAINER ID        IMAGE                      COMMAND             CREATED             STATUS              PORTS                    NAMES
666a5cc99d9b        company-b:latest           "/app-b"            5 minutes ago       Up 5 minutes        0.0.0.0:8082->8080/tcp   company-b
43617f7dd3a6        public-car-rental:latest   "/api"              5 minutes ago       Up 5 minutes        0.0.0.0:8083->8080/tcp   public-car-rental
39f64139e71f        company-a:latest           "/app-a"            5 minutes ago       Up 5 minutes        0.0.0.0:8081->8080/tcp   company-a
```

curl api
```bash
# Company A
curl -v -H "token:35216c9e-dea4-458c-babd-325f2ef0eefb" localhost:8081/data

# Company B
curl -v -H "token:35216c9e-dea4-458c-babd-325f2ef0eefb" localhost:8082/data

# Public Car Rental API
curl -v localhost:8083/data
```
check logs output and payload


### Test
```bash
# run test
make test

# check test cover
make test-cover-all
```

##### Notes
* check `Makefile` for details
* check `configs/config.go` for environment variables
