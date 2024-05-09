package main

// Yapı (struct) tanımları
type Item struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}