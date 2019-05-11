#!/usr/bin/with-contenv bash

echo "Setting up the Plex Media Server environment"
home="$(echo ~plex)"
export PLEX_MEDIA_SERVER_APPLICATION_SUPPORT_DIR="${PLEX_MEDIA_SERVER_APPLICATION_SUPPORT_DIR:-${home}/Library/Application Support}"
export PLEX_MEDIA_SERVER_HOME=/usr/lib/plexmediaserver
export PLEX_MEDIA_SERVER_MAX_PLUGIN_PROCS=6
export PLEX_MEDIA_SERVER_INFO_VENDOR=Docker
export PLEX_MEDIA_SERVER_INFO_DEVICE="Docker Container"
export PLEX_MEDIA_SERVER_INFO_MODEL=$(uname -m)
export PLEX_MEDIA_SERVER_INFO_PLATFORM_VERSION=$(uname -r)

if [ ! -d "${PLEX_MEDIA_SERVER_APPLICATION_SUPPORT_DIR}" ]; then
  /bin/mkdir -p "${PLEX_MEDIA_SERVER_APPLICATION_SUPPORT_DIR}"
  chown plex:plex "${PLEX_MEDIA_SERVER_APPLICATION_SUPPORT_DIR}"
fi

exec s6-setuidgid plex /bin/sh -c 'LD_LIBRARY_PATH=/usr/lib/plexmediaserver:/usr/lib/plexmediaserver/lib /usr/local/bin/plexcluster transcoder'