package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := GetDatabase()
	handleErr(err)

	bookHandler := NewBook(db)

	app := fiber.New()

	app.Get("/", hello)
	app.Get("/api/v1/books", bookHandler.ListBook)
	app.Post("/api/v1/books", bookHandler.CreateBook)
	app.Get("/api/v1/books/:id", bookHandler.GetBook)
	app.Patch("/api/v1/books/:id", bookHandler.UpdateBook)
	app.Delete("/api/v1/books/:id", bookHandler.DeleteBook)

	app.Listen(":3000")
}

func hello(c *fiber.Ctx) error {
	return c.Send([]byte("Hello"))
}

type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

type BookHandler struct {
	db *gorm.DB
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var book Book
	if err := c.BodyParser(&book); err != nil {
		return c.JSON(err)
	}

	if err := h.db.Create(&book).Error; err != nil {
		return c.JSON(err)
	}

	return c.JSON(book)
}

func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if err := h.db.First(&book, id).Error; err != nil {
		return c.JSON(err)
	}
	return c.JSON(book)
}

func (h *BookHandler) ListBook(c *fiber.Ctx) error {
	var books []Book
	if err := h.db.Find(&books).Error; err != nil {
		return c.JSON(err)
	}
	return c.JSON(books)
}

func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")

	var book Book
	if err := h.db.First(&book, id).Error; err != nil {
		return c.JSON(err)
	}

	if err := c.BodyParser(&book); err != nil {
		c.JSON(err)
	}

	if err := h.db.Updates(&book).Error; err != nil {
		c.JSON(err)
	}

	return c.JSON(book)
}

func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	if err := h.db.First(&book, id).Error; err != nil {
		return c.JSON(err)
	}

	if err := h.db.Delete(&book).Error; err != nil {
		return c.JSON(err)
	}

	return c.Send([]byte(fmt.Sprintf("book id [%s] delete successfully.", id)))
}

func NewBook(db *gorm.DB) *BookHandler { return &BookHandler{db} }

func GetDatabase() (*gorm.DB, error) {
	conn, err := gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	conn.AutoMigrate(&Book{})
	return conn, err
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
