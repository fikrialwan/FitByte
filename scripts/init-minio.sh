#!/bin/sh

# Wait for MinIO to be ready
echo "Waiting for MinIO to be ready..."
until wget -q --spider http://minio:9000/minio/health/live 2>/dev/null; do
    echo "MinIO is not ready yet. Waiting..."
    sleep 2
done

echo "MinIO is ready. Creating bucket..."

# Download and install MinIO client
wget -O /tmp/mc https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x /tmp/mc

# Configure MinIO client
/tmp/mc alias set minio http://minio:9000 ${MINIO_ROOT_USER} ${MINIO_ROOT_PASSWORD}

# Create bucket if it doesn't exist
/tmp/mc mb minio/${MINIO_BUCKET} --ignore-existing

# Set public read policy for the bucket (optional, for easier access)
/tmp/mc anonymous set public minio/${MINIO_BUCKET}

echo "MinIO bucket '${MINIO_BUCKET}' created successfully!"
