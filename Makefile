# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=plexcluster-go

run-docker:
	docker-compose up --build
build-docker:
	docker build -t plex-server -f docker/plex.Dockerfile .
	docker build -t plex-worker -f docker/worker.Dockerfile .
build-app:
	protoc -I=./ --go_out=plugins=grpc:./ ./plexcluster/plexcluster.proto
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
