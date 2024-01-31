package middleware

import "github.com/gofiber/fiber/v2"

type AuthMiddleware fiber.Handler

type RoleMiddleware func(roles ...string) fiber.Handler
