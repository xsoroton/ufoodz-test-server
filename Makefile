GO ?= go
GODEP ?= dep
GOIMPORTS ?= ${GOPATH}/bin/goimports
GOFMT ?= gofmt
GOLINT ?= ${GOPATH}/bin/golint

SERVICE_NAME=mobileapi
BUILDDATE=`date -u +%Y-%m-%d\ %H:%M`
GOOS=linux
GOARCH=amd64

LOCALHOST_IP ?= $(shell hostname --ip-address)

TEST_PACKAGE_PATHS := $(shell ${GO} list ./... | grep -Ev '/vendor/|/mocks|/model')
PACKAGE_PATHS := $(shell ${GO} list ./... | grep -v /vendor/)
GOSOURCES := $(shell find . -type f -name '*.go' -not -path './vendor/*')

GO_BUILD := GOOS=${GOOS} GOARCH=${GOARCH} $(GO) build -a -tags netgo -ldflags "-X \"main.build=${BUILD}\" -X \"main.buildDate=${BUILDDATE}\""

JWT_SECRET ?= SpjDimBfySs24H5QOErfH95XzN2sXmzVcrLigggWLJA

deps:
	$(GODEP) ensure -v

run:
	COMPANY_A_PORT=8081 $(GO) run private-company-a/main.go & \
	COMPANY_B_PORT=8082 $(GO) run private-company-b/main.go & \
	PORT=8083 COMPANY_A_PORT=8081 COMPANY_B_PORT=8082 $(GO) run public-car-rental/main.go

build: build-binary build-docker

build-binary:
	${GO_BUILD} -o ./build/app-a ./private-company-a/
	${GO_BUILD} -o ./build/app-b ./private-company-b/
	${GO_BUILD} -o ./build/api ./public-car-rental/

build-docker: build-binary
	docker build -t company-a:latest -f private-company-a/Dockerfile .
	docker build -t company-b:latest -f private-company-b/Dockerfile .
	docker build -t public-car-rental:latest -f public-car-rental/Dockerfile .

# TODO: change it to docker-compose
run-containers:
	docker network create mynet & \
	docker run -i --name company-a --rm --net=mynet -e COMPANY_A_PORT=8080 -p 8081:8080 company-a:latest & \
	docker run -i --name company-b --rm --net=mynet -e COMPANY_B_PORT=8080 -p 8082:8080 company-b:latest & \
	docker run -i --name public-car-rental --rm --net=mynet -p 8083:8080 \
	-e PORT=8080 \
	-e COMPANY_A_HOST=http://${LOCALHOST_IP} \
	-e COMPANY_A_PORT=8081 \
	-e COMPANY_B_HOST=http://${LOCALHOST_IP} \
	-e COMPANY_B_PORT=8082 \
	public-car-rental:latest

clean-up:
	docker rm -f company-a company-b public-car-rental

# run the tests
test:
	$(GO) test -v ${TEST_PACKAGE_PATHS}


# run the tests with coverage
test-cover:
	$(GO) test -covermode=count -coverprofile=coverage.out ${TEST_PACKAGE_PATHS}

# run tests with coverage and outputs coverage to console
#
# See https://blog.golang.org/cover for more output options
test-cover-stats: test-cover
	@if [ -e coverage.out ]; then \
		$(GO) tool cover -func=coverage.out; \
	fi
test-cover-html: test-cover
	@if [ -e coverage.out ]; then \
		$(GO) tool cover -html=coverage.out -o coverage.html; \
	fi
test-cover-all: test-cover-stats test-cover-html
test-cover-open: test-cover
	$(GO) tool cover -html=coverage.out

# run goimports validation on packages
imports:
	$(call validate,$(GOIMPORTS) -d,${GOSOURCES},GoImports,Lib)

# run gofmt with simplify validation on packages
fmt-simplify:
	$(call validate,$(GOFMT) -d -s,${GOSOURCES},GoFMT Simplify,Lib)

# run golint validation on packages
lint:
	$(call validate,$(GOLINT),${PACKAGE_PATHS},GoLint,Lib)

# run all validation on packages
validate: lint fmt-simplify imports

# code to run for all validation tests
define validate
	$(eval validate_test := `$(1) $(2) 2>&1`)
	@if [ "${validate_test}" = "" ]; then \
		echo "$(3) $(NAME) $(4): Good format"; \
	else \
		echo "$(3) $(NAME) $(4): Bad format"; \
		echo "${validate_test}"; \
		exit 33; \
	fi
endef
