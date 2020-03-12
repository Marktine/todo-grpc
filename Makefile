# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME= server
BINARY_UNIX=$(BINARY_NAME)_unix
CD_CMD_SERVER= cd ./cmd/server

all: test build
build:
	$(CD_CMD_SERVER) && \
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(CD_CMD_SERVER) && \
	$(GOTEST) -v 
clean:
	$(CD_CMD_SERVER) &&\
	$(GOCLEAN) &&\
	rm -f $(BINARY_NAME) &&\
	rm -f $(BINARY_UNIX) &&\
build-linux:
	$(CD_CMD_SERVER) &&\
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v