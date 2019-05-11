FROM golang:alpine as builder

RUN apk add --no-cache git
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -extldflags "-static"' -o plexcluster-go

FROM plexinc/pms-docker
ENV SERVER_URI=""
ENV TRANSCODER_PATH="/usr/lib/plexmediaserver/Plex Transcoder"
RUN touch /plex.log && \
    chown plex:plex /plex.log && \
    chmod 777 /plex.log
RUN cp "$TRANSCODER_PATH" /usr/lib/plexmediaserver/plex_transcoder
ADD docker/worker.sh /etc/services.d/plex/run
COPY --from=builder /app/plexcluster-go /usr/local/bin/plexcluster
RUN chmod +x /usr/local/bin/plexcluster