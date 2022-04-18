#!/bin/bash

# The entrypoint for this docker container. This script is responsible for
# configuring Tailscale within this container, and then starting the application.

echo "Starting tailscaled..."

tailscaled --socks5-server=localhost:1080 \
    --state=/run/tailscale-storage/tailscale.state \
    --tun=userspace-networking \
    --socket=/run/tailscale-storage/tailscale.sock &

echo "Started tailscaled."

# If we have $TAILSCALE_KEY (e.g., it isn't empty), then we authenticate.
if [ ! -z "$TAILSCALE_KEY" ]; then
    echo "Authenticating with Tailscale key..."
    until tailscale --socket=/run/tailscale-storage/tailscale.sock \
        up \
        --authkey=$TAILSCALE_KEY
    do 
        echo "waiting..."
        sleep 5
    done
    echo "Authenticated with Tailscale."
fi

# Sleep for ~5 seconds to give tailscaled time to start up.
sleep 5

# Run the application.
exec /run/server