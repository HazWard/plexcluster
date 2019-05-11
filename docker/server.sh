#!/usr/bin/with-contenv bash

exec s6-setuidgid plex /bin/sh -c '/usr/local/bin/plexcluster server'