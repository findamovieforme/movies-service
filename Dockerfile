# Stage 1: Build Go
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o movies-service main.go

# Stage 2: Runtime (Debian slim) - reliable sklearn wheels
FROM python:3.12-slim

# certificates for HTTPS calls
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Python deps
COPY helpers/requirements.txt /app/helpers/requirements.txt
RUN pip install --no-cache-dir -r /app/helpers/requirements.txt

# Copy Go binary, Python app, and model files
COPY --from=builder /app/movies-service /app/movies-service
COPY helpers/predictor.py /app/helpers/predictor.py
COPY helpers/app.py /app/helpers/app.py
COPY entrypoint.sh /app/entrypoint.sh
COPY recommendation-model /app/recommendation-model

RUN chmod +x /app/entrypoint.sh

EXPOSE 8081 5000
CMD ["/app/entrypoint.sh"]
