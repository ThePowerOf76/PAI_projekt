package middleware

import (
	"gdoc/model"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func CheckIfRegisteredToken(c *fiber.Ctx) error {
	token := c.Cookies("authtoken", "null")
	len, err := model.CheckTokenValidity(token)
	if err != nil {
		return c.RedirectToRoute("login", fiber.Map{"queries": map[string]string{"referred": "true"}})

	}
	if len > 0 {
		return c.Next()
	}
	return c.RedirectToRoute("login", fiber.Map{"queries": map[string]string{"referred": "true"}})
}
func JWTMiddlewareWrapper() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte("mySupah$ecretKeyL0LSneedChuckFuckAndSuck")},
		TokenLookup: "cookie:authtoken",
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return c.RedirectToRoute("login", fiber.Map{"queries": map[string]string{"referred": "true"}})
		},
	})
}
