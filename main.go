package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"www.github.com/techbrolakes/go-fiber-postgres/models"
	"www.github.com/techbrolakes/go-fiber-postgres/storage"
)

type Repository struct {
	DB *gorm.DB
}

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := new(Book)

	if err := c.BodyParser(book); err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "request failed"})
		return err
	}

	if err := r.DB.Create(book).Error; err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "Could not create book"})
		return err
	}

	if err := r.DB.Find(book).Error; err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"message": "Book Already exists"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book has been added", "data": book})
	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	if err := r.DB.Find(bookModels).Error; err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Could not get books"})
		return err
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books Fetched Successfully", "data": bookModels})
	return nil
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	id := c.Params("id")

	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Id Cannot be empty"})
		return nil
	}
	if err := r.DB.Delete(bookModels, id); err.Error != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not delete book"})
		return err.Error
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books Deleted Successfully"})
	return nil
}

func (r *Repository) GetBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	bookModels := &[]models.Books{}
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Id Cannot be empty"})
		return nil
	}

	fmt.Println("The Id is", id)
	if err := r.DB.Where("id = ?", id).First(bookModels).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get the book"})
		return err
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book Id Fetched Successfully", "data": bookModels})
	return nil
}
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not establish a database connection")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate DB")
	}

	r := &Repository{DB: db}
	app := fiber.New()
	r.SetupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
