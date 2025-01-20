package misc

import (
	"os"
	"vexal/structs"

	"github.com/gofiber/fiber/v3"
)

// GET /api/version
func GetVersion(c fiber.Ctx) error {
	return c.JSON(structs.Response{
		Error: false,
		Code:  200,
		Data: map[string]interface{}{
			"major":  os.Getenv("API_VERSION_MAJOR"),
			"minor":  os.Getenv("API_VERSION_MINOR"),
			"commit": os.Getenv("API_VERSION_COMMIT"),
		},
	})
}
