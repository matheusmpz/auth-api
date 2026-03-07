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
    var user models.LoginInput
    
    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais inválidas"})
        return
    }

    // Limpa email 
    cleanEmail := strings.ToLower(strings.TrimSpace(user.Email))

    //Busca usuário no banco
    var userID int
    var userName string
    var userEmail string
    var storedPassword string
    var isActive bool
    var isBlocked bool

    err := ctrl.DB.QueryRow(
		` SELECT id, name, email, password, is_active, is_blocked
          FROM users 
          WHERE email = $1`, 
		cleanEmail,
	).Scan(&userID, &userName, &userEmail, &storedPassword, &isActive, &isBlocked)

    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
        } else {
            log.Println("Erro ao buscar usuário:", err)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro no servidor"})
        }
        return
    }

    // Valida senha 
    err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
        return
    }

    // Valida se está bloqueado
    if isBlocked {
        ctx.JSON(http.StatusForbidden, gin.H{"error": "Conta bloqueada. Entre em contato com o suporte."})
        return
    }

    // Valida se está ativo
    if !isActive {
        ctx.JSON(http.StatusForbidden, gin.H{"error": "Conta não ativada. Verifique seu email."})
        return
    }

    // Retorna dados do usuário
    ctx.JSON(http.StatusOK, gin.H{
        "message": "Login realizado com sucesso",
        "user": gin.H{
            "id":    userID,
            "name":  userName,
            "email": userEmail,
        },
    })
}

func (ctrl *UserController) Activate(ctx *gin.Context) {
	var user models.ActivateInput

	// Validação de formato do email e código
	if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
        return
    }

    cleanEmail := strings.ToLower(strings.TrimSpace(user.Email))
    cleanCode := strings.TrimSpace(user.Code)

	// Fazendo a busca pelo código e se está ativo
	var activationCode sql.NullString
    var isActive bool
    var userID int

    err := ctrl.DB.QueryRow(`
        SELECT id, activation_code, is_active 
        FROM users 
        WHERE email = $1
    `, cleanEmail,).Scan(&userID, &activationCode, &isActive)

    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "E-mail ou código incorretos"})
        } else {
            log.Println("Erro ao buscar usuário:", err)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro no servidor"})
        }
        return
    }

	// Valida se já está ativo
    if isActive {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Conta já está ativada"})
        return
    }

    // Valida se código existe (não é NULL)
    if !activationCode.Valid {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Código de ativação não encontrado"})
        return
    }

    // Compara códigos
    if activationCode.String != cleanCode {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "E-mail ou código incorretos"})
        return
    }

    // Ativa a conta (UPDATE)
    _, err = ctrl.DB.Exec(`
        UPDATE users 
        SET is_active = true, activation_code = NULL
        WHERE id = $1
    `, userID)

    if err != nil {
        log.Println("Erro ao ativar conta:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ativar conta"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Conta ativada com sucesso!"})
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
