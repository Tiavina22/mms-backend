# MMS Backend

Complete backend for mobile messaging application built with Go, PostgreSQL and WebSocket.

## Features

- **JWT Authentication** - Secure signup/login
- **Encrypted Messages** - AES-256-GCM encryption
- **Group Messaging** - Create and manage groups
- **WebSocket** - Real-time communication
- **Push Notifications** - FCM & APNs support
- **Internationalization** - FR, EN, ES

## Quick Start

```bash
# 1. Install dependencies
go mod tidy

# 2. Configure environment
cp .env.sample .env
# Edit .env with your configuration

# 3. Run application
go run cmd/main.go

# 4. Run tests
go test ./tests/... -v
```

## Commands

```bash
go run cmd/main.go        # Start server
go test ./tests/... -v    # Run tests (like npm test)
go build cmd/main.go      # Build binary
make help                 # Show all commands
```

## Technology Stack

**Go** 1.21+ | **PostgreSQL** | **Gin** | **GORM** | **WebSocket**

## Documentation

- [SETUP.md](SETUP.md) - Detailed installation guide
- [COMMANDS.md](COMMANDS.md) - All available commands
- [API_EXAMPLES.md](API_EXAMPLES.md) - API usage examples
- [DOCS.md](DOCS.md) - Documentation index

## API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - Register new user
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/me` - Get current user

### Messages
- `POST /api/v1/messages` - Send message
- `GET /api/v1/messages/conversation/:id` - Get conversation
- `GET /api/v1/messages/conversations` - List conversations
- `PUT /api/v1/messages/read/:id` - Mark as read
- `GET /api/v1/messages/unread/count` - Unread count

### Groups
- `POST /api/v1/groups` - Create group
- `GET /api/v1/groups/my` - List my groups
- `POST /api/v1/groups/messages` - Send group message
- `GET /api/v1/groups/:id/messages` - Get group messages

### Users
- `GET /api/v1/users` - List users
- `GET /api/v1/users/search?q=term` - Search users
- `GET /api/v1/users/:id` - Get user details

### WebSocket
- `ws://localhost:8080/api/v1/ws` - Real-time connection

## Tests

**28 integration tests - 100% passing**

```bash
go test ./tests/... -v
```

**Coverage:**
- Authentication (signup, login, JWT)
- Direct messages (with encryption)
- Group messaging
- Security & authorization
- WebSocket communication

## Security

- **JWT** - Token-based authentication
- **AES-256-GCM** - Message encryption
- **Bcrypt** - Password hashing
- **Input validation** - All inputs sanitized
- **CORS** - Configured and secure

## Environment Configuration

Create `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=mms_db

PORT=8080
JWT_SECRET=your-secret-key
ENCRYPTION_KEY=32-byte-key-here

FCM_SERVER_KEY=your-fcm-key
APNS_KEY_ID=your-apns-key
```

See `.env.sample` for complete configuration.

## Project Structure

```
cmd/          # Application entry point
config/       # Configuration management
controllers/  # API endpoint handlers
models/       # Database models
repositories/ # Data access layer
services/     # Business logic
routes/       # Route definitions
utils/        # Utilities (JWT, crypto, validation)
websocket/    # WebSocket hub & clients
locales/      # i18n translations
tests/        # Integration tests
```

## Docker

```bash
docker build -t mms-backend .
docker run -p 8080:8080 --env-file .env mms-backend
```

## Usage Examples

### Signup

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"Test1234!"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"john@example.com","password":"Test1234!"}'
```

### Send Message

```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"receiver_id":"uuid","content":"Hello!"}'
```

## Development

### Run in development mode

```bash
ENV=development go run cmd/main.go
```

### Run tests

```bash
go test ./tests/... -v
make test
```

### Build for production

```bash
go build -o mms-backend cmd/main.go
make build
```

## Contributing

1. Fork the project
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details

## Built With

- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [Gorilla WebSocket](https://github.com/gorilla/websocket) - WebSocket implementation
- [golang-jwt](https://github.com/golang-jwt/jwt) - JWT implementation

---

**Production Ready** - Configure your secrets and deploy!
