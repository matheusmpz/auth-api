package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/matheusmpz/auth-api/models"
	"github.com/matheusmpz/auth-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	DB *sql.DB
}

var (
	validName = regexp.MustCompile(`^[a-zA-ZÀ-ÿ\s]+$`)
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func (ctrl *UserController) Register(ctx *gin.Context) {
    var user models.RegisterInput

    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
        return
    }

	// Preparando o dado para validação
    cleanName := strings.TrimSpace(user.Name)
    nameLength := utf8.RuneCountInString(cleanName)

	// Limpando o nome que depois será salvo no banco caso não retorne erro
	user.Name = cleanName

	// Validação de tamanho do nome
    if nameLength <= 2 || nameLength >= 100 {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "O nome deve ter entre 3 e 100 caracteres"})
        return
    }

	// Validação de formato valido do nome
    if !validName.MatchString(cleanName) {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "O nome deve conter apenas letras"})
        return
    }

	// Validação de email
	cleanEmail := strings.ToLower(strings.TrimSpace(user.Email))
    user.Email = cleanEmail
    if !emailRegex.MatchString(cleanEmail) {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "E-mail inválido"})
        return
    }

	// Validação de senha
	if pLen := utf8.RuneCountInString(user.Password); pLen < 6 || pLen > 255 {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "A senha deve ter no mínimo 6 caracteres"})
        return
    }

	// Verificação se o e-mail já não foi cadastrado
	var exists int
    err := ctrl.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&exists)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar email"})
        return
    }
    if exists > 0 {
        ctx.JSON(http.StatusConflict, gin.H{"error": "Email já cadastrado"})
        return
    }

	// Gerando o código de ativação
	activationCode := utils.GenerateActivationCode()

	// Colocando hash na senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
        return
    }

	// Salvando o usuário no banco de dados
	var userID int
    err = ctrl.DB.QueryRow(`
        INSERT INTO users (name, email, password, activation_code, is_active)
        VALUES ($1, $2, $3, $4, false)
        RETURNING id
    `, user.Name, user.Email, string(hashedPassword), activationCode).Scan(&userID)

	if err != nil {
		log.Println("Erro ao criar usuário:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

    // Retorna sucesso e o código de ativação para conseguir realizar 
	// o login sem precisar ativar a mensageria agora
    ctx.JSON(http.StatusCreated, gin.H{
        "message":         "Usuário criado com sucesso! Verifique seu email para ativar.",
        "user_id":         userID,
        "activation_code": activationCode,
    })
}

func (ctrl *UserController) Login(ctx *gin.Context) {
}

func (ctrl *UserController) Activate(ctx *gin.Context) {
}

func (ctrl *UserController) GetUser(ctx *gin.Context) {
}

func (ctrl *UserController) UpdateUser(ctx *gin.Context) {
}

func (ctrl *UserController) DeleteUser(ctx *gin.Context) {
}

func (ctrl *UserController) ActivateUser(ctx *gin.Context) {
}

func (ctrl *UserController) BlockUser(ctx *gin.Context) {
}

func (ctrl *UserController) UnblockUser(ctx *gin.Context) {
}
