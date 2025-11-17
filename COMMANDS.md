# 🚀 Commandes MMS Backend

Guide rapide des commandes disponibles (similaire à `npm test` en Node.js)

## 📦 Installation

```bash
# Installer les dépendances
go mod tidy

# ou avec Make
make deps
```

## 🏃 Lancer l'application

```bash
# Méthode 1: Directement avec Go
go run cmd/main.go

# Méthode 2: Avec Make
make run

# Méthode 3: Build puis exécuter
go build -o mms-backend cmd/main.go
./mms-backend
```

## 🧪 Tests (comme `npm test`)

```bash
# Lancer tous les tests (RECOMMANDÉ)
go test ./tests/... -v

# Version courte (sans verbose)
go test ./tests/...

# Avec timeout
go test ./tests/... -timeout 30s

# Avec Make
make test
```

## 📊 Coverage

```bash
# Générer le rapport de couverture
make coverage

# Ouvre automatiquement coverage.html dans le navigateur
```

## 🔧 Build

```bash
# Build pour votre OS
go build -o mms-backend cmd/main.go

# ou avec Make
make build

# Build pour Linux (depuis Windows/Mac)
GOOS=linux GOARCH=amd64 go build -o mms-backend cmd/main.go

# Build pour Windows (depuis Linux/Mac)
GOOS=windows GOARCH=amd64 go build -o mms-backend.exe cmd/main.go
```

## 🧹 Nettoyage

```bash
# Nettoyer les fichiers générés
make clean

# ou manuellement
go clean
rm -f mms-backend coverage.out coverage.html
```

## 📋 Aide

```bash
# Voir toutes les commandes Make disponibles
make help

# ou juste
make
```

---

## 🎯 Workflow de développement

### 1. **Setup initial**
```bash
go mod tidy
```

### 2. **Développement**
```bash
# Terminal 1: Lancer l'app
go run cmd/main.go

# Terminal 2: Tester
go test ./tests/... -v
```

### 3. **Avant de commit**
```bash
# Lancer les tests
go test ./tests/...

# Build pour vérifier
go build cmd/main.go
```

### 4. **Déploiement**
```bash
# Build production
go build -o mms-backend cmd/main.go

# ou avec optimisations
go build -ldflags="-s -w" -o mms-backend cmd/main.go
```

---

## 📝 Équivalences Node.js ↔ Go

| Node.js | Go | Description |
|---------|-------|-------------|
| `npm install` | `go mod tidy` | Installer les dépendances |
| `npm start` | `go run cmd/main.go` | Lancer l'application |
| `npm test` | `go test ./tests/...` | Lancer les tests |
| `npm run build` | `go build` | Compiler l'application |
| `npm run clean` | `go clean` | Nettoyer les fichiers |

---

## ✨ Raccourcis utiles

```bash
# Tout en un: deps + test + build
go mod tidy && go test ./tests/... && go build cmd/main.go

# Watch mode (nécessite air ou reflex)
# Install: go install github.com/cosmtrek/air@latest
air

# Format code
go fmt ./...

# Linter
go vet ./...
```

---

## 🎓 Commandes avancées

```bash
# Tests avec race detector
go test ./tests/... -race

# Tests avec coverage détaillée
go test ./tests/... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out

# Benchmark (si des tests de bench existent)
go test ./tests/... -bench=.

# Tests spécifiques
go test ./tests/... -run TestSignup
go test ./tests/... -run TestAuth

# Verbose avec temps
go test ./tests/... -v -timeout 30s
```

---

## 🐳 Docker (optionnel)

```bash
# Build image
docker build -t mms-backend .

# Run container
docker run -p 8080:8080 --env-file .env mms-backend
```

---

**Note**: Pas besoin de scripts PowerShell ou Bash complexes.  
Utilisez directement `go test` comme vous utiliseriez `npm test` ! 🚀

