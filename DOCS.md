# 📚 Documentation Index - MMS Backend

Index de toute la documentation disponible.

## 📖 Guides Principaux

| Fichier | Description |
|---------|-------------|
| [README.md](README.md) | **Démarrage rapide** - Vue d'ensemble du projet |
| [SETUP.md](SETUP.md) | **Installation détaillée** - Guide complet d'installation |
| [COMMANDS.md](COMMANDS.md) | **Commandes** - Toutes les commandes Go/Make disponibles |
| [API_EXAMPLES.md](API_EXAMPLES.md) | **Exemples API** - Exemples curl pour tous les endpoints |

---

## 🚀 Pour Commencer

**Nouveau sur le projet ?** Suivez cet ordre :

1. **[README.md](README.md)** → Vue d'ensemble et quick start
2. **[SETUP.md](SETUP.md)** → Installation complète
3. **[COMMANDS.md](COMMANDS.md)** → Commandes de développement
4. **[API_EXAMPLES.md](API_EXAMPLES.md)** → Tester l'API

---

## 📋 Quick Links

### Installation
- [Prerequisites](SETUP.md#prerequisites)
- [Database Setup](SETUP.md#step-3-create-database)
- [Environment Config](SETUP.md#step-4-configure-environment-variables)

### Development
- [Run Application](COMMANDS.md#-lancer-lapplication)
- [Run Tests](COMMANDS.md#-tests-comme-npm-test)
- [Build Production](COMMANDS.md#-build)

### API
- [Authentication](API_EXAMPLES.md#-authentication)
- [Messages](API_EXAMPLES.md#-direct-messages)
- [Groups](API_EXAMPLES.md#-groups)
- [WebSocket](API_EXAMPLES.md#-websocket)

---

## 🏗️ Architecture

```
mms-backend/
├── 📄 README.md           # Vue d'ensemble
├── 📘 SETUP.md            # Installation détaillée
├── 📗 COMMANDS.md         # Guide des commandes
├── 📙 API_EXAMPLES.md     # Exemples d'utilisation API
├── 📕 DOCS.md             # Ce fichier
│
├── cmd/                   # Point d'entrée
├── config/                # Configuration
├── controllers/           # Endpoints API
├── models/                # Modèles de données
├── repositories/          # Accès aux données
├── services/              # Logique métier
├── routes/                # Définition des routes
├── utils/                 # Utilitaires (JWT, crypto, etc)
├── websocket/             # Hub WebSocket
├── locales/               # Traductions i18n
└── tests/                 # Tests d'intégration
```

---

## 🎯 Par Besoin

### Je veux...

**...installer le projet**
→ [SETUP.md](SETUP.md)

**...lancer l'application**
```bash
go run cmd/main.go
```
→ [COMMANDS.md](COMMANDS.md#-lancer-lapplication)

**...lancer les tests**
```bash
go test ./tests/... -v
```
→ [COMMANDS.md](COMMANDS.md#-tests-comme-npm-test)

**...tester l'API**
→ [API_EXAMPLES.md](API_EXAMPLES.md)

**...comprendre l'architecture**
→ [README.md](README.md#-structure)

**...ajouter une fonctionnalité**
→ Voir la structure MVC dans `controllers/`, `services/`, `repositories/`

**...déployer en production**
→ [README.md](README.md#-docker) + Build avec `go build`

---

## 🧪 Tests

**28 tests d'intégration - 100% passent** ✅

```bash
go test ./tests/... -v
```

Tests disponibles dans `tests/integration_test.go`:
- ✅ Authentification (JWT)
- ✅ Messages directs (encryption AES-256)
- ✅ Groupes
- ✅ Sécurité
- ✅ WebSocket

---

## 🔧 Technologies

| Technologie | Usage | Doc |
|-------------|-------|-----|
| **Go** 1.21+ | Langage | [golang.org](https://golang.org) |
| **Gin** | Web Framework | [gin-gonic.com](https://gin-gonic.com) |
| **GORM** | ORM | [gorm.io](https://gorm.io) |
| **PostgreSQL** | Base de données | [postgresql.org](https://postgresql.org) |
| **JWT** | Auth | [jwt.io](https://jwt.io) |
| **WebSocket** | Temps réel | [gorilla/websocket](https://github.com/gorilla/websocket) |

---

## 📦 Structure des Packages

```
models/          → Structures de données (User, Message, Group, etc)
repositories/    → Accès DB (CRUD operations)
services/        → Logique métier (business logic)
controllers/     → Handlers HTTP (API endpoints)
routes/          → Route definitions
middleware/      → Auth, CORS, etc
utils/           → JWT, encryption, validation, i18n
websocket/       → Hub, clients, handlers
locales/         → Fichiers de traduction JSON
tests/           → Tests d'intégration
```

---

## 🔐 Sécurité

- **JWT** - Tokens avec expiration
- **AES-256-GCM** - Chiffrement des messages
- **Bcrypt** - Hash des mots de passe
- **Validation** - Inputs sanitizés
- **CORS** - Configuré et sécurisé

---

## 🌍 Multilingual (i18n)

Langues supportées : **FR**, **EN**, **ES**

Fichiers dans `locales/`:
- `en.json` - English
- `fr.json` - Français
- `es.json` - Español

Ajoutez une langue : créez `locales/de.json` (par exemple)

---

## 🚀 Déploiement

### Docker
```bash
docker build -t mms-backend .
docker run -p 8080:8080 --env-file .env mms-backend
```

### Build manuel
```bash
go build -o mms-backend cmd/main.go
./mms-backend
```

### Variables d'environnement
Voir `.env.sample` ou [SETUP.md](SETUP.md#step-4-configure-environment-variables)

---

## 🤝 Contribution

1. Fork le projet
2. Créer une branche (`git checkout -b feature/NewFeature`)
3. Commit (`git commit -m 'Add NewFeature'`)
4. Push (`git push origin feature/NewFeature`)
5. Pull Request

---

## 📞 Support

- **Issues** : [GitHub Issues](../../issues)
- **Documentation** : Ce dossier
- **Code** : Commenté en anglais

---

## 📝 License

MIT License - Voir LICENSE pour détails

---

**💡 Astuce**: Utilisez `Ctrl+F` pour rechercher rapidement dans cette documentation !

