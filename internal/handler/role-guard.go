package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maximfedotov74/fiber-psql/internal/utils"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

func (h *Handler) roleGuard(roles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userId, err := utils.GetUserIdFromCtx(ctx)

		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		user, err := h.services.UserService.GetUserById(*userId)

		if err != nil {
			return ctx.Status(err.Status()).JSON(err)
		}

		rolesFound := 0
		mustRolesFound := len(roles)

		for _, role := range roles {
			for _, userRole := range user.Roles {
				if userRole.Title == role {
					rolesFound++
					break
				}
			}
		}

		if rolesFound == mustRolesFound {
			return ctx.Next()
		}

		forbidden := lib.NewErr(messages.FORBIDDEN, 403)
		return ctx.Status(forbidden.Status()).JSON(forbidden)
	}
}
