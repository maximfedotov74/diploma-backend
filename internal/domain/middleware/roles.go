package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/utils"
)

// TODO add correct Role struct
func CreateRoleMiddleware() RoleMiddleware {
	return func(roles ...string) fiber.Handler {
		return func(ctx *fiber.Ctx) error {

			contextData, err := utils.GetLocalSession(ctx)

			if err != nil {
				return ctx.Status(err.Status()).JSON(err)
			}

			rolesFound := 0
			mustRolesFound := len(roles)

			for _, role := range roles {
				for _, userRole := range contextData.Roles {
					if strings.ToUpper(userRole.Title) == strings.ToUpper(role) {
						rolesFound++
						break
					}
				}
			}

			if rolesFound == mustRolesFound {
				return ctx.Next()
			}

			forbidden := fall.NewErr(fall.FORBIDDEN, fall.STATUS_FORBIDDEN)
			return ctx.Status(forbidden.Status()).JSON(forbidden)
		}
	}
}
