package routes

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	databaseconn "github.com/jiteshchawla1511/url-shortener/databaseConn"
	"github.com/jiteshchawla1511/url-shortener/helper"
	"github.com/jiteshchawla1511/url-shortener/models"
)

func Shorten(ctx *fiber.Ctx) error {
	body := &models.Request{}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	r2 := databaseconn.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(databaseconn.Ctx, ctx.IP()).Result()
	limit, _ := r2.TTL(databaseconn.Ctx, ctx.IP()).Result()

	if err != nil {
		_ = r2.Set(databaseconn.Ctx, ctx.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()

	} else if err == nil {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	if !govalidator.IsURL(body.URL) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	}
	if !helper.RemoveDomainError(body.URL) {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Kya reh bevde domain sahi daal"})
	}

	body.URL = helper.EnforceHTTP(body.URL)
	var id string
	if body.CustomShort == "" {
		id = helper.Base62Encode(rand.Uint64())
	} else {
		id = body.CustomShort
	}

	r1 := databaseconn.CreateClient(0)

	defer r1.Close()

	val, _ = r1.Get(databaseconn.Ctx, id).Result()
	if val != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL Custom short already in use",
		})
	}

	err = r1.Set(databaseconn.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to connect to server",
		})
	}

	defaultAPIQuotaStr := os.Getenv("API_QUOTA")

	defaultApiQuota, _ := strconv.Atoi(defaultAPIQuotaStr)
	resp := models.Response{
		URL:            body.URL,
		CustomShort:    "",
		Expiry:         body.Expiry,
		RateRemaining:  defaultApiQuota,
		RateLimitReset: 30,
	}

	remainingQuota, err := r2.Decr(databaseconn.Ctx, ctx.IP()).Result()

	resp.RateRemaining = int(remainingQuota)
	resp.RateRemaining = int(limit / time.Nanosecond / time.Minute)

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return ctx.Status(fiber.StatusOK).JSON(resp)

}
