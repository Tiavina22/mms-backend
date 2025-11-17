# Makefile pour MMS Backend

.PHONY: help run test coverage clean build deps

# Variables
APP_NAME=mms-backend
GO=go

help: ## Afficher cette aide
	@echo "Commandes disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Lancer l'application
	$(GO) run cmd/main.go

build: ## Compiler l'application
	$(GO) build -o $(APP_NAME) cmd/main.go

test: ## Lancer les tests
	$(GO) test -v ./tests/... -timeout 30s

coverage: ## Générer le rapport de couverture
	$(GO) test ./tests/... -coverprofile=coverage.out -timeout 30s
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Rapport de couverture: coverage.html"

clean: ## Nettoyer les fichiers générés
	rm -f $(APP_NAME) coverage.out coverage.html
	$(GO) clean

deps: ## Installer les dépendances
	$(GO) mod download
	$(GO) mod tidy

.DEFAULT_GOAL := help
