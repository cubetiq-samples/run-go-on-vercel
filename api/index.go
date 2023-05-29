package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Handler is the main entry point of the application. Think of it like the main() method
func Handler(w http.ResponseWriter, r *http.Request) {
	// This is needed to set the proper request path in `*fiber.Ctx`
	r.RequestURI = r.URL.String()

	handler().ServeHTTP(w, r)
}

// building the fiber application
func handler() http.HandlerFunc {
	app := fiber.New()

	// Optional middleware
	// app.Use("/ws", func(c *fiber.Ctx) error {
	// 	if c.Get("host") == "localhost:3000" {
	// 		c.Locals("Host", "Localhost:3000")
	// 		return c.Next()
	// 	}
	// 	return c.Status(403).SendString("Request origin not allowed")
	// })

	// Upgraded websocket request
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		fmt.Println(c.Locals("Host")) // "Localhost:3000"
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	app.Get("/v1", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"version": "v1",
		})
	})

	app.Get("/v2", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"version": "v2",
		})
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"uri":  ctx.Request().URI().String(),
			"path": ctx.Path(),
		})
	})

	return adaptor.FiberApp(app)
}
