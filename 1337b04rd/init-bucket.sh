#!/bin/sh

set -e

echo "⏳ Waiting for MinIO to be available..."
sleep 5

mc alias set local http://minio:9000 "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD"

if ! mc ls local/board-storage > /dev/null 2>&1; then
  echo "✅ Creating bucket 'board-storage'..."
  mc mb local/board-storage
else
  echo "✅ Bucket 'board-storage' already exists"
fi
