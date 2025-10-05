package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	// Serve embedded static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal("Failed to create static filesystem:", err)
	}

	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(staticFS),
	}))

	// Deploy endpoint
	app.Post("/deploy/:serviceName", deployHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func deployHandler(c *fiber.Ctx) error {
	serviceName := c.Params("serviceName")
	apiKey := c.Get("X-API-Key")

	// Validate API key
	expectedKey := os.Getenv("NATTUKAKA_API_KEY")
	if expectedKey == "" {
		return c.Status(500).JSON(fiber.Map{
			"error": "NATTUKAKA_API_KEY not configured",
		})
	}

	if apiKey != expectedKey {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid API key",
		})
	}

	// Check if service exists
	deployScript := filepath.Join("/root/nattukaka/apps", serviceName, "deploy.sh")
	if _, err := os.Stat(deployScript); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("Service '%s' not found", serviceName),
		})
	}

	// Execute deploy script and capture output
	cmd := exec.Command("bash", deployScript)
	cmd.Dir = filepath.Join("/root/nattukaka/apps", serviceName)

	// Combine stdout and stderr
	output, err := cmd.CombinedOutput()

	if err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		return c.Status(500).SendString(fmt.Sprintf("Deploy failed with exit code %d:\n%s", exitCode, string(output)))
	}

	return c.SendString(fmt.Sprintf("Deployment completed successfully:\n%s", string(output)))
}
