package models

import "time"

// Struct do usuário definida relacionando corretamente~
// com os campos do banco de dados
// realizado um pequeno tratamento dos campos de ID e Password
// para que o ID só seja exibido apenas se não for vazio (omitempty)
// para o senha e codigo de ativação foi anulado a exibição no json (-)
// no código de ativação foi criado como o tipo ponteiro de string pois pode vir nulo
type User struct {
	ID             int       `json:"id,omitempty" db:"id"`
	Name           string    `json:"name" db:"name"`
	Email          string    `json:"email" db:"email"`
	Password       string    `json:"-" db:"password"`
	Active         bool      `json:"is_active" db:"is_active"`
	Blocked        bool      `json:"is_blocked" db:"is_blocked"`
	ActivationCode *string   `json:"-" db:"activation_code"`
	Created        time.Time `json:"created_at" db:"created_at"`
	Updated        time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterInput - Dados de cadastro
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput - Dados de login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ActivateInput - Dados de ativação
type ActivateInput struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// UserResponse - Resposta
type UserResponse struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Active  bool      `json:"is_active"`
	Blocked bool      `json:"is_blocked"`
	Created time.Time `json:"created_at"`
}