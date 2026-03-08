# Auth API - Sistema de Autenticação em Go

Sistema completo de autenticação desenvolvido em Go com PostgreSQL.

## 🚀 Features

- [ ] Cadastro de usuários
- [ ] Login com JWT
- [ ] Hash de senha (bcrypt)
- [ ] Ativação de conta
- [ ] Bloqueio/desbloqueio de conta

## 🛠️ Tecnologias

- Go 1.26
- PostgreSQL 14+
- Gin Framework
- bcrypt
- JWT

## 📦 Setup
```bash
# 1. Clone o repositório
git clone https://github.com/matheusmpz/auth-api.git
cd auth-api

# 2. Configure o .env
cp .env.example .env
# Edite o .env com suas credenciais

# 3. Crie o banco de dados
psql -U postgres -f scripts/setup.sql

# 4. Instale as dependências
go mod download

# 5. Rode a aplicação
go run cmd/api/main.go
```

## 🧪 Endpoints

### Rotas Públicas (sem autenticação)

| Método | Endpoint | Descrição | Body |
|--------|----------|-----------|------|
| POST | `/register` | Cadastra novo usuário | `{name, email, password}` |
| POST | `/login` | Autentica usuário e retorna JWT | `{email, password}` |
| POST | `/activate` | Ativa conta com código | `{email, code}` |

### Rotas Protegidas (requer JWT no header)

| Método | Endpoint | Descrição | Header |
|--------|----------|-----------|--------|
| GET | `/users/:id` | Busca usuário por ID | `Authorization: Bearer <token>` |
| PUT | `/users/:id` | Atualiza dados do usuário | `Authorization: Bearer <token>` |
| DELETE | `/users/:id` | Deleta usuário | `Authorization: Bearer <token>` |
| PATCH | `/users/:id/activate` | Ativa conta manualmente | `Authorization: Bearer <token>` |
| PATCH | `/users/:id/block` | Bloqueia conta | `Authorization: Bearer <token>` |
| PATCH | `/users/:id/unblock` | Desbloqueia conta | `Authorization: Bearer <token>` |
