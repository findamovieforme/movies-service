#!/bin/sh
set -e

# Start the Flask/Gunicorn recommendation service
gunicorn -w 2 -b 0.0.0.0:5000 helpers.app:app &

# Start the Go movies-service (foreground)
exec /app/movies-service

