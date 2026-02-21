#!/bin/sh

echo "Waiting for PostgreSQL to be ready..."

# Простая задержка вместо сложной проверки
sleep 5

exec ./main