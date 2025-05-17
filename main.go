package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type Task struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

//var tasks = map[uint]Task{}

func main() {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}
	log.Println("Database connected")

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Task{})

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/tasks", func(c *fiber.Ctx) error {
		var task Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		if err := db.Create(&task).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
		}
		return c.Status(fiber.StatusCreated).JSON(task)
	})

	app.Get("/tasks", func(c *fiber.Ctx) error {
		var taskList []Task
		if err := db.Find(&taskList).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve tasks"})
		}
		return c.JSON(taskList)
	})

	app.Get("/tasks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var task Task
		if err := db.First(&task, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		return c.JSON(task)
	})

	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var task Task
		if err := db.First(&task, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		var updatedTask Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}
		task.Title = updatedTask.Title
		task.Done = updatedTask.Done
		if err := db.Save(&task).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
		}
		return c.JSON(task)
	})

	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var task Task
		if err := db.First(&task, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		if err := db.Delete(&task).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete task"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	})

	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
