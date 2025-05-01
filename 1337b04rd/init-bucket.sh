#!/bin/sh

set -e

echo "⏳ Waiting for MinIO to be available..."
until curl -s http://minio:9000 >/dev/null; do
  echo "MinIO is not yet available, waiting..."
  sleep 5
done

mc alias set local http://minio:9000 "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD"

if ! mc ls local/board-storage > /dev/null 2>&1; then
  echo "✅ Creating bucket 'board-storage'..."
  mc mb local/board-storage
else
  echo "✅ Bucket 'board-storage' already exists"
fi
