# MinIO Setup Guide

This guide explains how to set up and use MinIO for file storage in the FitByte project.

## Overview

MinIO is an S3-compatible object storage server that provides:
- **Local Development**: No need for AWS credentials during development
- **Cost Effective**: Free for self-hosted deployments
- **S3 Compatible**: Uses the same AWS SDK, making migration seamless
- **Performance**: High-performance object storage
- **Web Console**: Built-in web interface for bucket management

## Configuration

### Environment Variables

Add these variables to your `.env` file:

```bash
# MinIO Configuration
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=fitbyte-uploads
MINIO_USE_SSL=false
MINIO_PORT=9000
MINIO_CONSOLE_PORT=9001
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
```

### Docker Compose Setup

MinIO is automatically configured in `docker-compose.yml` with:
- **MinIO Server**: Runs on port 9000 (API) and 9001 (Console)
- **Automatic Bucket Creation**: The `minio-init` service creates the required bucket
- **Persistent Storage**: Data is stored in `fit_byte_minio_data` volume
- **Health Checks**: Ensures MinIO is ready before starting the application

## Usage

### Starting the Services

```bash
# Start all services including MinIO
docker-compose up -d

# Check MinIO status
docker-compose ps minio
```

### Accessing MinIO Console

1. Open your browser and go to: http://localhost:9001
2. Login with:
   - **Username**: minioadmin
   - **Password**: minioadmin
3. You can browse buckets, upload files, and manage permissions

### File Upload API

The file upload endpoint remains the same:

```bash
POST /v1/file
Content-Type: multipart/form-data

# Upload a file
curl -X POST \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/image.jpg" \
  http://localhost:8080/v1/file
```

**Response:**
```json
{
  "uri": "http://localhost:9000/fitbyte-uploads/uploads/2024/01/15/unique-filename.jpg"
}
```

## Development vs Production

### Development (Local)
- Uses `localhost:9000` endpoint
- HTTP (no SSL)
- Default credentials (minioadmin/minioadmin)

### Production
Update these environment variables for production:
```bash
MINIO_ENDPOINT=your-minio-server.com:9000
MINIO_ACCESS_KEY=your-production-access-key
MINIO_SECRET_KEY=your-production-secret-key
MINIO_USE_SSL=true
MINIO_ROOT_USER=your-admin-user
MINIO_ROOT_PASSWORD=your-secure-password
```

## File Organization

Files are organized with the following structure:
```
fitbyte-uploads/
└── uploads/
    └── YYYY/
        └── MM/
            └── DD/
                └── unique-filename.ext
```

Example: `uploads/2024/01/15/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg`

## Troubleshooting

### Common Issues

1. **Bucket Not Found Error**
   - Ensure the `minio-init` service completed successfully
   - Check if the bucket exists in MinIO console

2. **Connection Refused**
   - Verify MinIO container is running: `docker-compose ps minio`
   - Check MinIO health: `curl http://localhost:9000/minio/health/live`

3. **Access Denied**
   - Verify credentials in environment variables
   - Check bucket permissions in MinIO console

### Logs

```bash
# View MinIO logs
docker-compose logs minio

# View bucket initialization logs
docker-compose logs minio-init

# View application logs
docker-compose logs app
```

## Migration from AWS S3

The migration from AWS S3 to MinIO is seamless because:
1. **Same API**: MinIO implements the S3 API
2. **Same SDK**: Uses the same AWS SDK for Go
3. **Same Code**: Only configuration changes, no code changes needed

### What Changed
- **Endpoint**: Points to MinIO server instead of AWS
- **Credentials**: Uses MinIO access keys instead of AWS IAM
- **URL Format**: Returns MinIO URLs instead of S3 URLs

### What Stayed the Same
- **File Upload Logic**: Identical multipart upload handling
- **Validation**: Same file size and type restrictions
- **API Interface**: Same REST endpoints and responses
- **Error Handling**: Same error responses and status codes

## Security Notes

- Change default credentials in production
- Use HTTPS in production (`MINIO_USE_SSL=true`)
- Configure proper bucket policies
- Consider network security (VPC, firewall rules)
- Regular backups of MinIO data volume
