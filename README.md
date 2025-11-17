# MMS Backend 💬

Backend complet pour application de messagerie mobile avec Go, PostgreSQL et WebSocket.

## ✨ Features

- 🔐 **Auth JWT** - Signup/Login sécurisé
- 💬 **Messages chiffrés** - AES-256-GCM
- 👥 **Groupes** - Messagerie de groupe
- 🌐 **WebSocket** - Temps réel
- 🔔 **Push** - FCM & APNs
- 🌍 **i18n** - FR, EN, ES

## 🚀 Quick Start

```bash
# 1. Install
go mod tidy

# 2. Configure
cp .env.sample .env
# Edit .env with your config

# 3. Run
go run cmd/main.go

# 4. Test
go test ./tests/... -v
```

## 📋 Commandes

```bash
go run cmd/main.go        # Lancer
go test ./tests/... -v    # Tester (comme npm test)
go build cmd/main.go      # Build
make help                 # Voir toutes les commandes
```

## 🔧 Stack

- **Go** 1.21+ | **PostgreSQL** | **Gin** | **GORM** | **WebSocket**

## 📚 Documentation

- [SETUP.md](SETUP.md) - Guide d'installation détaillé
- [COMMANDS.md](COMMANDS.md) - Toutes les commandes disponibles
- [API Endpoints](#api-endpoints) - Liste complète ci-dessous

## 🌐 API Endpoints

### Auth
- `POST /api/v1/auth/signup` - Inscription
- `POST /api/v1/auth/login` - Connexion
- `GET /api/v1/auth/me` - Utilisateur actuel

### Messages
- `POST /api/v1/messages` - Envoyer
- `GET /api/v1/messages/conversation/:id` - Conversation
- `GET /api/v1/messages/conversations` - Liste
- `PUT /api/v1/messages/read/:id` - Marquer lu
- `GET /api/v1/messages/unread/count` - Compteur

### Groups
- `POST /api/v1/groups` - Créer
- `GET /api/v1/groups/my` - Mes groupes
- `POST /api/v1/groups/messages` - Envoyer message
- `GET /api/v1/groups/:id/messages` - Messages

### Users
- `GET /api/v1/users` - Lister
- `GET /api/v1/users/search?q=term` - Rechercher
- `GET /api/v1/users/:id` - Détails

### WebSocket
- `ws://localhost:8080/api/v1/ws` - Connexion temps réel

## 🧪 Tests

**28 tests d'intégration - 100%** ✅

```bash
go test ./tests/... -v
```

**Tests couverts:**
- ✓ Auth (signup, login, JWT)
- ✓ Messages directs (encryption)
- ✓ Groupes
- ✓ Sécurité
- ✓ WebSocket

## 🔒 Sécurité

- **JWT** - Authentification
- **AES-256-GCM** - Chiffrement messages
- **Bcrypt** - Hash passwords
- **Validation** - Inputs
- **CORS** - Configuré

## 📝 Environment (.env)

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

Voir `.env.sample` pour la configuration complète.

## 📂 Structure

```
cmd/          # Entry point
config/       # Configuration
controllers/  # API endpoints
models/       # Database models
repositories/ # Data access
services/     # Business logic
routes/       # Route definitions
utils/        # Utilities (JWT, crypto, etc)
websocket/    # WebSocket hub & clients
locales/      # Translations (i18n)
tests/        # Integration tests
```

## 🐳 Docker

```bash
docker build -t mms-backend .
docker run -p 8080:8080 --env-file .env mms-backend
```

## 📖 Exemples d'utilisation

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

## 🤝 Contributing

1. Fork le projet
2. Créer une branche (`git checkout -b feature/AmazingFeature`)
3. Commit (`git commit -m 'Add AmazingFeature'`)
4. Push (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

## 📄 License

MIT License - voir [LICENSE](LICENSE) pour plus de détails

## 🙏 Built with

- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [GORM](https://gorm.io/) - ORM
- [Gorilla WebSocket](https://github.com/gorilla/websocket) - WebSocket
- [golang-jwt](https://github.com/golang-jwt/jwt) - JWT

---

**🚀 Production Ready** - Configurez vos secrets et c'est prêt !
