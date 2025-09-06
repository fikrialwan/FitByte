#!/bin/sh

echo "Starting MinIO initialization..."

# Wait for MinIO to be ready using curl (more reliable than wget)
echo "Waiting for MinIO to be ready..."
max_attempts=30
attempt=0
until curl -f http://minio:9000/minio/health/live >/dev/null 2>&1; do
    attempt=$((attempt + 1))
    if [ $attempt -ge $max_attempts ]; then
        echo "ERROR: MinIO failed to become ready after $max_attempts attempts"
        exit 1
    fi
    echo "MinIO is not ready yet. Waiting... (attempt $attempt/$max_attempts)"
    sleep 5
done

echo "MinIO is ready. Creating bucket using API calls..."

# Create bucket using MinIO REST API instead of downloading mc client
MINIO_URL="http://minio:9000"
BUCKET_NAME="${MINIO_BUCKET:-fitbyte-uploads}"

echo "Creating bucket: $BUCKET_NAME"

# Create bucket using PUT request
response=$(curl -s -w "%{http_code}" -X PUT \
    -H "Host: $BUCKET_NAME.minio:9000" \
    --user "${MINIO_ROOT_USER}:${MINIO_ROOT_PASSWORD}" \
    "$MINIO_URL/$BUCKET_NAME" \
    -o /tmp/response.txt)

if [ "$response" = "200" ] || [ "$response" = "409" ]; then
    echo "Bucket '$BUCKET_NAME' created successfully or already exists"
else
    echo "Failed to create bucket. HTTP response: $response"
    cat /tmp/response.txt
    # Don't exit with error if bucket already exists
    if [ "$response" != "409" ]; then
        exit 1
    fi
fi

# Set bucket policy to public read using MinIO API
echo "Setting bucket policy to public read..."
policy='{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"AWS": "*"},
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::'$BUCKET_NAME'/*"]
    }
  ]
}'

curl -s -X PUT \
    -H "Content-Type: application/json" \
    --user "${MINIO_ROOT_USER}:${MINIO_ROOT_PASSWORD}" \
    --data "$policy" \
    "$MINIO_URL/minio/admin/v3/set-bucket-policy?bucket=$BUCKET_NAME" \
    >/dev/null 2>&1

echo "MinIO bucket '$BUCKET_NAME' initialized successfully!"
echo "You can access MinIO console at: http://localhost:9001"
echo "Username: ${MINIO_ROOT_USER}"
echo "Password: ${MINIO_ROOT_PASSWORD}"
