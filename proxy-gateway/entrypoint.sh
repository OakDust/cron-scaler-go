#!/bin/sh

echo "Waiting for scale-handler to be ready..."

# Ждем пока scale-handler gRPC сервер будет готов
while ! nc -z scale-handler 50051; do
  sleep 1
done

echo "scale-handler is ready!"

echo "Starting proxy-gateway service..."
exec ./main