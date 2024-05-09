package main

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "log"
)

var pool *pgxpool.Pool

// Veritabanı bağlantısını kuran fonksiyon
func initializeDatabaseConnection() {
    dbUrl := "postgresql://myuser:mypassword@10.151.231.133:5432/mydatabase"
    var err error
    pool, err = pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        log.Fatalf("Veritabanına bağlanılamadı: %v", err)
    }
}