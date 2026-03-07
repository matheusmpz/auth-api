package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateActivationCode gera código de 6 dígitos para ativação do usuário
func GenerateActivationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}