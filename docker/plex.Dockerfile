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
RUN rm "$TRANSCODER_PATH"
ADD docker/server.sh /etc/services.d/plexcluster/run
ADD docker/ffmpeg.sh /usr/lib/plexmediaserver/
RUN mv /usr/lib/plexmediaserver/ffmpeg.sh "$TRANSCODER_PATH" && \
    chmod +x "$TRANSCODER_PATH"
COPY --from=builder /app/plexcluster-go /usr/local/bin/plexcluster
RUN chmod +x /usr/local/bin/plexcluster
