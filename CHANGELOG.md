# Changelog

## [1.1.0] - 2026-03-08

### Adicionado
- AUPI-07 - JWT authentication implementation
  - Token generation on login
  - Token validation middleware
  - Protected routes with JWT
  - User authentication system

### Modificado
- Login now returns JWT token along with user data
- All user routes now require authentication

### Segurança
- JWT-based authentication for protected routes
- Authorization header validation
- Token expiration (24 hours)

## [1.0.1] - 2026-03-08

### Adicionado
- AUPI-01 - Project structure initialization
- AUPI-02 - Initiating connection to database
- AUPI-03 - Adding route and user registration processing
- AUPI-04 - Adding account activation processing
- AUPI-05 - Adding login route processing
- AUPI-06 - Adding the version file

### Segurança
- Hash de senhas com bcrypt
- Validação de dados de entrada
- Prevenção de SQL injection