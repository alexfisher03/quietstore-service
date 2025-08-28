package handlers

import (
	"context"
	"time"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/db"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

func ReadyCheck(dbConn db.Pingable, s3c *s3.Client, bucket string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		if err := dbConn.Ping(ctx); err != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, "db not ready: "+err.Error())
		}

		_, err := s3c.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: &bucket,
		})
		if err != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, "s3 not ready: "+err.Error())
		}

		return c.JSON(fiber.Map{
			"status":  "ready",
			"service": "QuietStore",
		})
	}
}
