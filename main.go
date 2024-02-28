package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Task struct {
	ID int `json:"id"`
}






func main() {
	
	// Create a new Redis client 
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	// Create a new Fiber application
	app := fiber.New()

	// Check if the Redis client is connected
	app.Use(func(c *fiber.Ctx) error {
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			c.Status(fiber.StatusInternalServerError).SendString("Error connecting to Redis")
			return err
		}
		return c.Next()
	})



	// Create a new task using the POST method and push it to the Redis queue
	app.Post("/tasks", func (c *fiber.Ctx) error {

		// create a new task 
		task := new(Task)
	
		var err error // Declare the "err" variable

		if err := c.BodyParser(task); err != nil {
			c.Status(fiber.StatusBadRequest).SendString("Invalid request")
			return err
		}

		var taskJSON []byte
		// Marshal the task struct to JSON
		taskJSON, err = json.Marshal(task)

		if err != nil {
			c.Status(fiber.StatusInternalServerError).SendString("Marshal error")
			return err
		}

		// Pushing the task in the redis list using LPUSH command
		err = client.LPush(context.Background(), "tasks", taskJSON).Err()

		if err != nil {
			c.Status(fiber.StatusInternalServerError).SendString("LPUSH error")
			return err
		}

		// Return the task as a JSON response
		c.JSON("Task added successfully")
		return nil
	})


	//Routing to get the next task from the Redis queue

	app.Get("/tasks", func(c *fiber.Ctx) error {

		// Pop the next task from the Redis list using the LPOP command
		taskJson, err := client.RPop(context.Background(), "tasks").Result()

		if err != nil {
			c.Status(fiber.StatusInternalServerError).SendString("LPOP error")
			return err
		}

		// Unmarshal the task JSON to a task struct
		task := new(Task)
		err = json.Unmarshal([]byte(taskJson), task)

		if err != nil {
			c.Status(fiber.StatusInternalServerError).SendString("Unmarshal error")
			return err
		}

		// Return the task as a JSON response
		c.JSON(task)
		return nil
	})

	// Start the Fiber server on port 3000
	app.Listen(":3000")
	fmt.Println("Server started on port 3000")
	
}