package routes

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	databaseconn "github.com/jiteshchawla1511/url-shortener/databaseConn"
)

func Resolve(ctx *fiber.Ctx) error {
	url := ctx.Params("url")

	r := databaseconn.CreateClient(0)
	defer r.Close()

	value, err := r.Get(databaseconn.Ctx, url).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short-url not found in db"})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal error"})
	}

	rInr := databaseconn.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(databaseconn.Ctx, "counter")

	return ctx.Redirect(value, 301)
}
