package main

import (
	controller "gdoc/controller"
	middleware "gdoc/middleware"
	model "gdoc/model"
	view "gdoc/view"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

func main() {
	templating_engine := html.New("./view/views", ".html")
	app := fiber.New(fiber.Config{
		Views: templating_engine,
	})
	app.Use(cors.New())
	app.Static("/static", "./static")
	err := model.DBInit()
	if err != nil {
		panic(err)
	}
	app.Get("/register", view.RegisterRouter)
	app.Post("/register", controller.RegisterHandler)
	app.Get("/", view.MainPageView).Name("root")
	app.Get("/login", view.LoginRouter).Name("login")
	app.Post("/login", controller.LoginHandler)
	app.Get("/scoreboard", view.ScoreboardPageView)
	app.Get("/scores", controller.GetScores)
	app.Get("/logout", controller.LogoutHandler)

	auth := app.Group("/l", middleware.JWTMiddlewareWrapper())
	game := auth.Group("", middleware.CheckIfRegisteredToken)

	game.Get("/profilelist", controller.ProfileList)
	game.Get("/profile/:id", controller.GetProfileById)
	game.Post("/profile", controller.AddNewProfile)
	game.Patch("/profile/:id", controller.UpdateProfile)
	game.Delete("/profile/:id", controller.DeleteProfile)
	game.Get("/getname", controller.GetAuthenticatedUserData)
	game.Post("/scores", controller.AddNewScore)
	game.Get("/customise", view.CustomisePageView)
	game.Get("/game", view.GamePageView)
	app.Listen(":3000")
}
