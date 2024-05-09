package main

import (
    "github.com/gofiber/fiber/v2"
)

// GET kök rotası
func getRoot(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
}

// Tüm öğeleri veritabanından alıp JSON formatında döndüren fonksiyon
func getAllItems(c *fiber.Ctx) error {
    rows, err := pool.Query(c.Context(), "SELECT id, name, price FROM items")
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    defer rows.Close()

    var items []Item
    for rows.Next() {
        var item Item
        rows.Scan(&item.ID, &item.Name, &item.Price)
        items = append(items, item)
    }

    return c.JSON(items)
}

// Yeni bir öğeyi ekleyen POST rotası
func addItem(c *fiber.Ctx) error {
    var newItem Item
    if err := c.BodyParser(&newItem); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    _, err := pool.Exec(c.Context(), "INSERT INTO items (name, price) VALUES ($1, $2)", newItem.Name, newItem.Price)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusCreated).JSON(newItem)
}