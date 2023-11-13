package guards

import (
	"github.com/gofiber/fiber/v2"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
	"github.com/maximfedotov74/fiber-psql/internal/shared/utils"
)

type RoleGuard struct{}

func NewRoleGuard() *RoleGuard {
	return &RoleGuard{}
}

func (rg *RoleGuard) CheckRoles(roles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		contextData, err := utils.GetUserDataFromCtx(ctx)

		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		rolesFound := 0
		mustRolesFound := len(roles)

		for _, role := range roles {
			for _, userRole := range contextData.Roles {
				if userRole.Title == role {
					rolesFound++
					break
				}
			}
		}

		if rolesFound == mustRolesFound {
			return ctx.Next()
		}

		forbidden := exception.NewErr(messages.FORBIDDEN, 403)
		return ctx.Status(forbidden.Status()).JSON(forbidden)
	}
}
