# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
BINARY_NAME=mybinary
BINARY_UNIX=$(BINARY_NAME)_unix

all: run
build: 
		$(GOBUILD) -o $(BINARY_NAME) -v
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		DISPLAY=:9 $(GORUN) src/wm.go src/window.go src/socket.go
deps:
		$(GOGET) github.com/markbates/goth
		$(GOGET) github.com/markbates/pop