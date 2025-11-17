# MMS Backend - Detailed Setup Guide

This guide will walk you through setting up the MMS Backend from scratch.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Step-by-Step Setup](#step-by-step-setup)
3. [Running the Application](#running-the-application)
4. [Testing the API](#testing-the-api)
5. [WebSocket Testing](#websocket-testing)
6. [Common Issues](#common-issues)

## Prerequisites

### Required Software

1. **Go** (version 1.21 or higher)
   - Download from: https://golang.org/dl/
   - Verify installation: `go version`

2. **PostgreSQL** (version 12 or higher)
   - Download from: https://www.postgresql.org/download/
   - Verify installation: `psql --version`

3. **Git**
   - Download from: https://git-scm.com/downloads
   - Verify installation: `git --version`

### Optional Tools

- **Postman** or **Insomnia** for API testing
- **TablePlus** or **pgAdmin** for database management
- **wscat** for WebSocket testing: `npm install -g wscat`

## Step-by-Step Setup

### Step 1: Install Go

**Windows:**
```powershell
# Download installer from https://golang.org/dl/
# Run the installer and follow the prompts
# Verify installation
go version
```

**macOS:**
```bash
brew install go
go version
```

**Linux:**
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

### Step 2: Install PostgreSQL

**Windows:**
1. Download installer from https://www.postgresql.org/download/windows/
2. Run installer and remember the password you set
3. Add PostgreSQL bin directory to PATH

**macOS:**
```bash
brew install postgresql
brew services start postgresql
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### Step 3: Create Database

```bash
# Connect to PostgreSQL
psql -U postgres

# In psql shell, create database
CREATE DATABASE mms_db;

# Create a user (optional)
CREATE USER mms_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE mms_db TO mms_user;

# Exit psql
\q
```

### Step 4: Clone and Setup Project

```bash
# Navigate to your workspace
cd /path/to/your/workspace

# If you have the code, navigate to it
cd mms-backend

# Install dependencies
go mod tidy
```

### Step 5: Configure Environment Variables

Create a `.env` file in the root directory:

```bash
# Copy the template
cp ENV_TEMPLATE.txt .env

# Edit .env with your preferred editor
nano .env  # or vim, code, notepad++, etc.
```

**Minimum Required Configuration:**

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=mms_db
DB_SSLMODE=disable

# Server
PORT=8080
ENV=development

# JWT Secret (generate a random string)
JWT_SECRET=my-super-secret-jwt-key-for-development-only

# Encryption Key (must be 32 characters)
ENCRYPTION_KEY=12345678901234567890123456789012

# Push Notifications (optional for development)
FCM_SERVER_KEY=
APNS_KEY_ID=
```

**Generate Secure Secrets:**

```bash
# Generate JWT Secret (macOS/Linux)
openssl rand -base64 32

# Generate Encryption Key (macOS/Linux)
openssl rand -hex 32

# Windows PowerShell
Add-Type -AssemblyName System.Web
[System.Web.Security.Membership]::GeneratePassword(32,5)
```

### Step 6: Verify Installation

Check that all dependencies are installed:

```bash
go mod verify
go mod download
```

## Running the Application

### Development Mode

```bash
# From the project root
go run cmd/main.go
```

You should see output like:

```
Configuration loaded successfully
Database connection established successfully
Database migrations completed
Loaded translations for languages: [en es fr]
WebSocket hub started
Starting MMS Backend server on port 8080
Environment: development
WebSocket endpoint: ws://localhost:8080/api/v1/ws
```

### Build and Run

```bash
# Build the binary
go build -o mms-backend cmd/main.go

# Run the binary
./mms-backend  # Linux/macOS
mms-backend.exe  # Windows
```

### Running with Custom Port

```bash
# Set PORT in .env or export it
export PORT=9000
go run cmd/main.go
```

## Testing the API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "MMS Backend is running"
}
```

### 2. User Signup

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test1234!",
    "phone": "+1234567890",
    "language": "en"
  }'
```

Expected response:
```json
{
  "message": "Account created successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "username": "testuser",
      "email": "test@example.com",
      ...
    }
  }
}
```

### 3. User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "test@example.com",
    "password": "Test1234!"
  }'
```

### 4. Get Current User (Authenticated)

```bash
# Replace <YOUR_TOKEN> with the token from signup/login
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

### 5. Search Users

```bash
curl http://localhost:8080/api/v1/users/search?q=test&limit=10 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

### 6. Send Message

```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "receiver_id": "receiver-user-uuid",
    "content": "Hello, this is a test message!"
  }'
```

## WebSocket Testing

### Using wscat

```bash
# Install wscat if not already installed
npm install -g wscat

# Connect to WebSocket
wscat -c "ws://localhost:8080/api/v1/ws" \
  -H "Authorization: Bearer <YOUR_TOKEN>"

# Once connected, send a message
{"type": "new_message", "receiver_id": "uuid", "content": "Hello via WebSocket!"}

# Send ping
{"type": "ping"}
```

### Using JavaScript

```javascript
const token = "YOUR_JWT_TOKEN";
const ws = new WebSocket("ws://localhost:8080/api/v1/ws");

// Set authorization header (if your WebSocket client supports it)
// Or pass token in URL: ws://localhost:8080/api/v1/ws?token=${token}

ws.onopen = () => {
  console.log("Connected to WebSocket");
  
  // Send a message
  ws.send(JSON.stringify({
    type: "new_message",
    receiver_id: "uuid-of-receiver",
    content: "Hello!"
  }));
};

ws.onmessage = (event) => {
  console.log("Received:", JSON.parse(event.data));
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};
```

## Common Issues

### Issue 1: Database Connection Failed

**Error:** `Failed to connect to database`

**Solutions:**
- Verify PostgreSQL is running: `pg_isready`
- Check database credentials in `.env`
- Ensure database exists: `psql -U postgres -l`
- Check firewall settings

### Issue 2: Port Already in Use

**Error:** `bind: address already in use`

**Solutions:**
```bash
# Find process using port 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill the process or change PORT in .env
```

### Issue 3: JWT Token Invalid

**Error:** `invalid or expired token`

**Solutions:**
- Ensure you're including the full token with "Bearer " prefix
- Check token hasn't expired (default: 24h)
- Verify JWT_SECRET matches between signup/login and verification

### Issue 4: Module Not Found

**Error:** `cannot find module`

**Solutions:**
```bash
# Clean and reinstall dependencies
go clean -modcache
go mod tidy
go mod download
```

### Issue 5: Migration Errors

**Error:** `migration failed`

**Solutions:**
```bash
# Drop and recreate database
psql -U postgres -c "DROP DATABASE mms_db;"
psql -U postgres -c "CREATE DATABASE mms_db;"

# Run application again
go run cmd/main.go
```

### Issue 6: WebSocket Connection Refused

**Solutions:**
- Ensure you're using the correct WebSocket protocol (ws:// not wss://)
- Check that JWT token is valid and not expired
- Verify CORS settings if connecting from browser
- Check that WebSocket endpoint is `/api/v1/ws`

## Next Steps

1. **Explore the API**: Try all endpoints documented in README.md
2. **Test WebSocket**: Send real-time messages between users
3. **Add More Users**: Create multiple test accounts
4. **Create Groups**: Test group messaging functionality
5. **Configure Push Notifications**: Set up FCM and APNs
6. **Deploy**: Follow deployment guide for production

## Getting Help

If you encounter issues not covered here:

1. Check the main README.md
2. Review application logs
3. Check PostgreSQL logs
4. Enable debug mode: Set `ENV=development` in `.env`
5. Create an issue in the repository

## Development Tips

### Hot Reload (Optional)

Install air for hot reloading during development:

```bash
go install github.com/cosmtrek/air@latest
air
```

### Database GUI

Use TablePlus or pgAdmin to view database:
- Host: localhost
- Port: 5432
- Database: mms_db
- User: postgres
- Password: your_password

### API Testing Collection

Import the following into Postman:
- Base URL: http://localhost:8080/api/v1
- Create environment variables for token and user_id

### Debugging

Enable verbose logging:
```bash
# Set GIN_MODE
export GIN_MODE=debug
go run cmd/main.go
```

---

**Congratulations!** You now have a fully functional MMS Backend running locally.

