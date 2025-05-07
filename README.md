# HouseBank

HouseBank is a backend service for managing user accounts, transactions, and authentication. It is built using modern Go practices and integrates with various AWS services for deployment and management.

The project uses SQLC for database query generation, AWS ECR for container image storage, AWS EKS for Kubernetes-based deployment, AWS RDS for database management, and AWS Secrets Manager for secure storage of sensitive information.

## DB docs:
- [Db Docs](https://dbdocs.io/ankitnayan83/HouseBank)
- [Api Docs](https://app.swaggerhub.com/apis-docs/personal-05c/House_Bank_Api/1.0)

## Features

- **User Management**: Create, authenticate, and manage users
- **Account Management**: Create, update, and delete accounts with support for multiple currencies
- **Transaction Management**: Transfer money between accounts with concurrency-safe operations
- **Authentication**: Secure token-based authentication using JWT and PASETO
- **Database**: PostgreSQL database with schema migrations managed by golang-migrate
- **Testing**: Comprehensive unit tests using testify and gomock
- **Deployment**: Fully containerized application deployed on AWS infrastructure

## Technologies Used

### Backend
- **Go**: The application is written in Go (v1.24)
- **Gin**: A lightweight web framework for building RESTful APIs
- **SQLC**: Generates type-safe Go code from SQL queries
- **PostgreSQL**: Relational database for storing user and account data
- **bcrypt**: For secure password hashing

### AWS Services
- **AWS ECR**: Stores Docker images for deployment
- **AWS EKS**: Manages Kubernetes clusters for running the application
- **AWS RDS**: Hosts the PostgreSQL database
- **AWS Secrets Manager**: Securely stores environment variables and sensitive data

### DevOps
- **Docker**: Containerizes the application for consistent deployment
- **Kubernetes**: Orchestrates containerized workloads
- **GitHub Actions**: Automates CI/CD pipelines for testing and deployment
- **golang-migrate**: Manages database schema migrations

## Project Structure

```
├── api/                 # API handlers and middleware
├── db/                  # Database migrations and queries
│   ├── migration/       # SQL migration files
│   ├── query/           # SQL query files
│   └── sqlc/            # Generated Go code
├── util/                # Utility functions
├── config/              # Configuration management
├── token/               # Authentication token management
├── kubernetes/          # Kubernetes deployment files
├── workflows/           # GitHub Actions workflow files
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── Makefile             # Build automation
├── app.env.example      # Example environment variables
└── main.go              # Application entry point
```

## Setup and Installation

### Prerequisites
- Go (v1.24 or later)
- Docker
- AWS CLI
- PostgreSQL

### Steps

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/housebank.git
   cd housebank
   ```

2. **Set Up Environment Variables**
   Create an `app.env` file with the following variables:
   ```
   DB_DRIVER=postgres
   DB_SOURCE=postgresql://root:password@localhost:5432/housebank?sslmode=disable
   SERVER_ADDRESS=0.0.0.0:8080
   TOKEN_SYMMETRIC_KEY=your_symmetric_key_at_least_32_characters
   ACCESS_TOKEN_DURATION=15m
   ```

3. **Run Database Migrations**
   ```bash
   make migrateup
   ```

4. **Run the Application**
   ```bash
   make server
   ```

5. **Run Tests**
   ```bash
   make test
   ```

## Deployment

### Docker

1. **Build Docker Image**
   ```bash
   make image
   ```

2. **Run Docker Container**
   ```bash
   docker run -p 8080:8080 -e DB_SOURCE=your_db_source housebank:latest
   ```

### AWS Deployment

1. **Push Docker Image to AWS ECR**
   ```bash
   aws ecr get-login-password --region YOUR_REGION | docker login --username AWS --password-stdin YOUR_ACCOUNT_ID.dkr.ecr.YOUR_REGION.amazonaws.com
   docker tag housebank:latest YOUR_ACCOUNT_ID.dkr.ecr.YOUR_REGION.amazonaws.com/housebank:latest
   docker push YOUR_ACCOUNT_ID.dkr.ecr.YOUR_REGION.amazonaws.com/housebank:latest
   ```

2. **Deploy to AWS EKS**
   Use Kubernetes manifests or Helm charts to deploy the application to an EKS cluster.

3. **Configure AWS RDS**
   Set up a PostgreSQL instance on AWS RDS and update the `DB_SOURCE` in `app.env`.

4. **Use AWS Secrets Manager**
   Store sensitive environment variables in AWS Secrets Manager and load them during deployment.

## CI/CD

The project uses GitHub Actions for CI/CD. The workflows are defined in the `workflows` directory.

### Test Workflow
- Runs on every push or pull request to the main branch
- Sets up a PostgreSQL service and runs all tests

### Deploy Workflow
- Runs on every push to the main branch
- Builds a Docker image, pushes it to AWS ECR, and deploys it to AWS EKS

## Key Commands

### Makefile Commands
- Run PostgreSQL: `make postgresrun`
- Start Server: `make server`
- Run Tests: `make test`
- Build Docker Image: `make image`
- Run Migrations: `make migrateup`

### Docker Compose
- Start all services: `docker-compose up`

## Testing

The project includes comprehensive unit tests for all major components:
- **Database Tests**: Tests for SQLC-generated queries
- **API Tests**: Tests for API endpoints using httptest
- **Mocking**: Uses gomock to mock database interactions

Run all tests:
```bash
make test
```

## Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Submit a pull request to the main branch

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contact

For any questions or issues, please contact Ankit Nayan.
