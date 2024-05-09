package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Veritabanı bağlantısını başlat
    initializeDatabaseConnection()
    defer pool.Close()

    // Fiber uygulamasını başlat
    app := fiber.New()

    // Rotaları tanımla
    app.Get("/", getRoot)
    app.Get("/items", getAllItems)
    app.Post("/items", addItem)

    // Sunucuyu başlat
    log.Fatal(app.Listen(":3000"))
}