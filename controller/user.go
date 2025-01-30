package controller

import (
	b64 "encoding/base64"
	"gdoc/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
	"time"
)

type LoginDetails struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func LoginHandler(c *fiber.Ctx) error {
	login_details := new(LoginDetails)
	if err := c.BodyParser(login_details); err != nil {
		return err
	}
	hash := argon2.IDKey([]byte(login_details.Password), []byte(login_details.Username), 1, 64*1024, 4, 32)
	userid, password, err := model.GetUserCredentials(login_details.Username)
	if err != nil {
		return err
	}
	if userid == 0 || password != b64.StdEncoding.WithPadding(-1).EncodeToString(hash) {
		return c.RedirectToRoute("login", fiber.Map{"queries": map[string]string{"err": "t"}})
	}
	expiry_date := time.Now().Add(time.Hour * 24)
	tokenClaims := jwt.MapClaims{
		"name":    login_details.Username,
		"uid":     userid,
		"expires": expiry_date.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	t, err := token.SignedString([]byte("mySupah$ecretKeyL0LSneedChuckFuckAndSuck"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = model.AddNewToken(t, expiry_date.Format(time.DateTime))
	if err != nil {
		return err
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "authtoken"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)
	return c.RedirectToRoute("root", fiber.Map{})
}
func LogoutHandler(c *fiber.Ctx) error {
	token := c.Cookies("authtoken", "null")
	if token != "null" {
		model.InvalidateToken(token)
	}
	c.ClearCookie("authtoken")
	return c.Render("logout", fiber.Map{})
}
func RegisterHandler(c *fiber.Ctx) error {
	login_details := new(LoginDetails)
	if err := c.BodyParser(login_details); err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return nil
	}

	r, err := model.GetUserCountByName(login_details.Username)
	if err != nil {
		return err
	}
	if r == 0 {
		hash := argon2.IDKey([]byte(login_details.Password), []byte(login_details.Username), 1, 64*1024, 4, 32)

		err = model.AddNewUser(login_details.Username, b64.StdEncoding.WithPadding(-1).EncodeToString(hash))
		if err != nil {
			c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   err,
			})
			return nil
		}
		return c.RedirectToRoute("login", fiber.Map{})
	} else {
		c.Status(409).RedirectToRoute("register", fiber.Map{"queries": map[string]string{"taken": "true"}})
		return nil
	}
}
func GetAuthenticatedUserData(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	uid := int(claims["uid"].(float64))
	return c.JSON(fiber.Map{
		"name": name,
		"uid":  uid,
	})
}
