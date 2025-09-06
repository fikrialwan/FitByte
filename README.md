# FitByte 🏃‍♂️💪

A modern fitness tracking application built with Go, providing a comprehensive REST API for managing user fitness activities and health metrics.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Usage Examples](#usage-examples)
- [Development](#development)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

FitByte is a RESTful API service designed to help users track their fitness activities, monitor their health metrics, and manage their fitness journey. The application provides secure user authentication, activity logging, file uploads, and comprehensive fitness data management.

## ✨ Features

### 🔐 User Management

- User registration and authentication
- JWT-based security
- Profile management with customizable preferences
- Support for different measurement units (weight/height)
- Profile image upload

### 🏃 Activity Tracking

- Support for 10+ activity types:
  - **Low Intensity**: Walking, Yoga, Stretching (4 cal/min)
  - **Medium Intensity**: Cycling, Swimming, Dancing (8 cal/min)
  - **High Intensity**: Hiking, Running, HIIT, Jump Rope (10 cal/min)
- Automatic calorie calculation based on activity duration
- Activity history and analytics
- Customizable activity preferences

### 📁 File Management

- Secure file upload to AWS S3
- Profile image management
- File validation and processing

### 🚀 Performance & Security

- Redis-based caching
- Rate limiting middleware
- API documentation with Swagger
- Comprehensive error handling
- Database migrations

## 🛠 Tech Stack

- **Backend**: Go 1.24
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **File Storage**: AWS S3
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or later
- PostgreSQL 16+
- Redis 7+
- AWS Account (for S3 file storage)
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/fikrialwan/FitByte.git
   cd FitByte
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

### Configuration

Configure your `.env` file with the following variables:

```env
# Database Configuration
DB_USER=your_db_user
DB_PASS=your_db_password
DB_HOST=localhost
DB_NAME=fitbyte
DB_PORT=5432

# Application Configuration
JWT_SECRET=your_super_secret_jwt_key
APP_PORT=8080
APP_ENV=develop

# AWS S3 Configuration
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_REGION=us-east-1
AWS_S3_BUCKET=your_s3_bucket_name

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
```

## 📚 API Documentation

The API documentation is automatically generated using Swagger and is available at:

```
http://localhost:8080/swagger/index.html
```

### Available Endpoints

- **Authentication**: `POST /v1/register`, `POST /v1/login`
- **User Management**: `GET /v1/user`, `PATCH /v1/user`
- **Activity Tracking**: `POST /v1/activity`, `GET /v1/activity`, `PATCH /v1/activity/:activityId`, `DELETE /v1/activity`
- **File Upload**: `POST /v1/file`

## 💻 Usage Examples

### User Registration

```bash
curl -X POST http://localhost:8080/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Log Activity

```bash
curl -X POST http://localhost:8080/v1/activity \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "activityType": "Running",
    "durationInMinutes": 30,
    "doneAt": "2024-01-15T08:00:00Z"
  }'
```

## 🔧 Development

### Available Make Commands

```bash
# Build the application
make build

# Run the application
make run

# Build and run
make build-run

# Run tests
make test

# Run database migrations
make migrate

# Generate Swagger documentation
make docs
```

### Running with Docker

1. **Start all services**

   ```bash
   docker-compose up -d
   ```

2. **View logs**

   ```bash
   docker-compose logs -f app
   ```

3. **Stop services**
   ```bash
   docker-compose down
   ```

### Project Structure

```
├── cmd/
│   ├── app/          # Main application
│   └── migrate/      # Database migrations
├── config/           # Configuration files
├── internal/
│   ├── controller/   # HTTP handlers
│   ├── dto/          # Data transfer objects
│   ├── entity/       # Database models
│   ├── repository/   # Data access layer
│   ├── routes/       # Route definitions
│   └── service/      # Business logic
├── middlewares/      # HTTP middlewares
├── docs/            # Swagger documentation
```

## 🚀 Deployment

The application can be deployed using Docker:

1. **Build production image**

   ```bash
   docker build -t fitbyte:latest .
   ```

2. **Deploy with Docker Compose**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### 🚀 Deploy to Kubernetes

Please read the [deployment documentation](deployments/k8s/README.md) to prepare and deploy this application to Kubernetes.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](http://www.apache.org/licenses/LICENSE-2.0.html) file for details.

---

**Happy Coding! 🚀**
