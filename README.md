# CipherSafe

A self-hosted secrets manager built with Go and Next.js. Securely store, encrypt, and manage your API keys, passwords, and sensitive data.

## Features

- **End-to-End Encryption**: All secrets are encrypted using AES-256-GCM
- **Project Organization**: Group secrets by projects for better management
- **Self-Hosted**: Keep full control over your data and infrastructure
- **Modern UI**: Clean, responsive interface built with Next.js and Tailwind CSS
- **JWT Authentication**: Secure authentication with JSON Web Tokens
- **PostgreSQL Database**: Robust data storage with GORM ORM

## Tech Stack

- **Backend**: Go, Gin, GORM, PostgreSQL
- **Frontend**: Next.js, React, TypeScript, Tailwind CSS
- **Authentication**: JWT tokens
- **Encryption**: AES-256-GCM

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- PostgreSQL database
- npm or yarn

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/haripatel07/CipherSafe.git
cd CipherSafe
```

### 2. Backend Setup

```bash
cd backend

# Install Go dependencies
go mod tidy

# Copy and configure environment variables
cp .env.example .env
# Edit .env with your database credentials and encryption keys

# Run the backend server
go run main.go
```

The backend will start on `http://localhost:8080`

### 3. Frontend Setup

```bash
cd ../frontend

# Install dependencies
npm install

# Run the development server
npm run dev
```

The frontend will be available at `http://localhost:3000`

### 4. Database Setup

Create a PostgreSQL database named `ciphersafe` and update the `DATABASE_URL` in your `.env` file.

The application will automatically create the necessary tables on first run.

## Environment Configuration

Create a `.env` file in the `backend` directory with the following variables:

```env
# A 32-byte (256-bit) key for AES encryption
MASTER_ENCRYPTION_KEY="your-32-byte-base64-encoded-key"

# A secret for signing JWT tokens
JWT_SECRET_KEY="your-super-secret-jwt-key"

# Database connection string
DATABASE_URL="host=localhost user=postgres password=yourpassword dbname=ciphersafe port=5432 sslmode=disable"
```

### Generating Encryption Keys

To generate a secure encryption key:

```bash
openssl rand -base64 32
```

## API Documentation

### Authentication Endpoints

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login and receive JWT token

### Protected Endpoints (require Bearer token)

- `POST /api/projects` - Create a new project
- `GET /api/projects` - Get all projects for the authenticated user
- `POST /api/secrets` - Create a new secret
- `GET /api/projects/:projectID/secrets` - Get all secrets for a project
- `DELETE /api/secrets/:secretID` - Delete a secret

## Usage

1. **Register**: Create an account at `/register`
2. **Login**: Sign in at `/login`
3. **Create Project**: Add a new project to organize your secrets
4. **Add Secrets**: Store encrypted key-value pairs within projects
5. **Manage**: View, copy, and delete secrets as needed

## Security Features

- **Encryption**: All secret values are encrypted before database storage
- **Authentication**: JWT-based authentication with secure token handling
- **Authorization**: Users can only access their own projects and secrets
- **HTTPS Ready**: Designed to work with HTTPS in production

## Development

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests (if implemented)
cd ../frontend
npm test
```

### Building for Production

```bash
# Backend
cd backend
go build -o ciphersafe main.go

# Frontend
cd ../frontend
npm run build
npm start
```

## Deployment

### Docker (Recommended)

Create a `Dockerfile` and `docker-compose.yml` for easy deployment:

```yaml
# docker-compose.yml
version: '3.8'
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: ciphersafe
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - db
    environment:
      - DATABASE_URL=host=db user=postgres password=password dbname=ciphersafe port=5432 sslmode=disable

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
```

### Manual Deployment

1. Set up PostgreSQL database
2. Configure environment variables
3. Build and run the backend
4. Build and serve the frontend (Vercel, Netlify, etc.)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

If you encounter any issues or have questions:

1. Check the existing issues on GitHub
2. Create a new issue with detailed information
3. Include your environment setup and error messages

## Roadmap

- [ ] Multi-user project sharing
- [ ] Secret versioning
- [ ] CLI tool for secret management
- [ ] Integration with popular CI/CD platforms
- [ ] Audit logging
- [ ] Backup and restore functionality
