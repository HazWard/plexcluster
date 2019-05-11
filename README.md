# PlexCluster

A distributed transcoding infrastructure for Plex using gRPC


## Dependencies

* docker-compose 1.23.1+

For development, add these as well:

* Go 1.11+
* protoc 3.0+

## Usage

For now the project requires more testing and documentation.

To run a Plex server and a worker node, simply run:
```
docker-compose up --build
```