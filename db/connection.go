package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "sync"
    "time"

    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

var (
    db   *sql.DB // ponteiro da conexão com o banco, só é preenchida dentro do onde.Do
    once sync.Once // controla que o bloco once.Do rode apenas uma vez
)

func GetDB() *sql.DB {
    once.Do(func() {
        // Carregando o arquivo .env
        err := godotenv.Load()
        if err != nil {
            log.Println("Variáveis de ambiente incorretas")
        }

        // Monta a string de conexão com as variáveis do .env
        connStr := fmt.Sprintf(
            "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
            os.Getenv("DB_HOST"),
            os.Getenv("DB_PORT"),
            os.Getenv("DB_USER"),
            os.Getenv("DB_PASSWORD"),
            os.Getenv("DB_NAME"),
            os.Getenv("DB_SSLMODE"),
        )

        // prepara a conexão com o banco, o sql.Open não faz a conexão
        // ele apenas válida o formato da string e prepara o driver
        conn, err := sql.Open("postgres", connStr)
        if err != nil {
            log.Fatal("Erro ao abrir conexão:", err)
        }

        // Aqui temos um gerenciamento de conexões
        // O `*sql.DB` não é uma conexão única — é um pool que gerencia várias conexões simultâneas. 
        // Sem configurar isso, ele abre conexões ilimitadas e pode derrubar o Postgres.

        conn.SetMaxOpenConns(25) // maximo de conexões simultâneas, caso passe de 25 as conexões ficam aguardando
        conn.SetMaxIdleConns(10) // conexões que podem ser reaproveitadas
        conn.SetConnMaxLifetime(5 * time.Minute) // tempo máximo de uma conexão, após isso ela finaliza e recria 
        conn.SetConnMaxIdleTime(2 * time.Minute) // tempo máximo de uma conexão no pool, se ninguém utilizar finaliza

        // Realiza a real conexão com o banco
        if err := conn.Ping(); err != nil {
            log.Fatal("Erro de conexão no banco:", err)
        }

        db = conn
    })

    // retorna a conexão para as chamadas
    return db
}