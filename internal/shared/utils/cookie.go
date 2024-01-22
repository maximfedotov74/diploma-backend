package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/shared/jwt"
)

func SetCookies(tokens jwt.Tokens) (*fiber.Cookie, *fiber.Cookie) {

	access_cookie := new(fiber.Cookie)
	access_cookie.Name = "access_token"
	access_cookie.Value = tokens.AccessToken
	access_cookie.Expires = tokens.AccessExpTime
	refresh_cookie := new(fiber.Cookie)
	refresh_cookie.Name = "refresh_token"
	refresh_cookie.Value = tokens.RefreshToken
	refresh_cookie.Expires = tokens.RefreshExpTime
	refresh_cookie.HTTPOnly = true

	return access_cookie, refresh_cookie
}
