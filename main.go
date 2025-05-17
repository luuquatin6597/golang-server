package main

import "github.com/gofiber/fiber/v2"

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = map[int]Task{}
var nextID = 1

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/tasks", func(c *fiber.Ctx) error {
		var task Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		task.ID = nextID
		nextID++
		tasks[task.ID] = task
		return c.Status(fiber.StatusCreated).JSON(task)
	})

	app.Get("/tasks", func(c *fiber.Ctx) error {
		var taskList []Task
		for _, task := range tasks {
			taskList = append(taskList, task)
		}
		return c.JSON(taskList)
	})

	app.Get("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
		}
		task, exists := tasks[id]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		return c.JSON(task)
	})

	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
		}
		var updatedTask Task
		if err := c.BodyParser(&updatedTask); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}
		task, exists := tasks[id]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		task.Title = updatedTask.Title
		task.Done = updatedTask.Done
		tasks[id] = task
		return c.JSON(task)
	})

	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
		}
		_, exists := tasks[id]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		delete(tasks, id)
		return c.SendStatus(fiber.StatusNoContent)
	})

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
