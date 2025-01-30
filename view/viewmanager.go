package view

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRouter(c *fiber.Ctx) error {
	if c.Query("taken", "false") == "true" {
		return c.Render("register", fiber.Map{"RedirectBanner": "This username is already registered."})
	} else {
		return c.Render("register", fiber.Map{})
	}
}
func LoginRouter(c *fiber.Ctx) error {
	if c.Query("referred", "false") == "true" {
		return c.Render("login", fiber.Map{"RedirectBanner": "You must log in to continue."})
	} else if c.Query("err", "") == "t" {
		return c.Render("login", fiber.Map{"RedirectBanner": "Incorrect login or password"})
	} else {
		return c.Render("login", fiber.Map{})
	}
}
func MainPageView(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}
func CustomisePageView(c *fiber.Ctx) error {
	return c.Render("customise", fiber.Map{})
}
func GamePageView(c *fiber.Ctx) error {
	return c.Render("gdoc", fiber.Map{})
}
func ScoreboardPageView(c *fiber.Ctx) error {
	return c.Render("scoreboard", fiber.Map{})
}
