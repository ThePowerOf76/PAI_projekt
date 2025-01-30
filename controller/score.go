package controller

import (
	"gdoc/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Score struct {
	PID   float64 `json:"pid"`
	Score float64 `json:"score"`
}
type ScoreDisplay struct {
	Cursor string  `json:"cursor"`
	Colour string  `json:"colour"`
	Name   string  `json:"name"`
	Score  float64 `json:"score"`
}

func GetScores(c *fiber.Ctx) error {
	GetRow, err := model.GetScores()
	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}
	var name string
	var cursor string
	var color string
	var score float64
	var scores []ScoreDisplay
	for {
		cursor, color, name, score, err = GetRow()
		if err != nil {
			c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   err,
			})
			return err
		}
		if score == -1 {
			break
		}
		scores = append(scores, ScoreDisplay{Name: name, Cursor: cursor, Colour: color, Score: score})

	}
	c.Status(200).JSON(fiber.Map{"scores": scores})
	return nil
}
func AddNewScore(c *fiber.Ctx) error {
	c.Accepts("application/json")
	score := new(Score)
	if err := c.BodyParser(score); err != nil {
		return err
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	player_id := claims["uid"].(float64)

	err := model.AddNewScore(player_id, score.PID, score.Score)
	if err != nil {
		return err
	}
	return nil
}
