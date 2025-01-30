package controller

import (
	"gdoc/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"strconv"
)

type Profile struct {
	PID       float64 `json:"pid"`
	OID       float64 `json:"oid"`
	Name      string  `json:"name"`
	Cursor    string  `json:"cursor"`
	BGType    string  `json:"bgtype"`
	BGContent string  `json:"bgcontent"`
	PColour   string  `json:"pcolour"`
	SColour   string  `json:"scolour"`
	Segments  float64 `json:"segments"`
	Music     bool    `json:"music"`
	SFX       bool    `json:"sfx"`
	Active    bool    `json:"active"`
}

func ProfileList(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := int(claims["uid"].(float64))
	GetRow, err := model.GetProfiles(uid)
	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}
	var oid float64
	var pid float64
	var name string
	var cursor string
	var bgType string
	var bgContent string
	var pcolour string
	var scolour string
	var segments float64
	var music bool
	var sfx bool
	var active bool
	var profiles []Profile
	for {
		pid, oid, name, cursor, bgType, bgContent, pcolour, scolour, segments, music, sfx, active, err = GetRow()
		if err != nil {
			c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   err,
			})
			return err
		}
		if pid == -1 {
			break
		}
		if !active {
			continue
		}
		profiles = append(profiles, Profile{PID: pid, OID: oid, Name: name, Cursor: cursor, BGType: bgType, BGContent: bgContent, PColour: pcolour, SColour: scolour, Segments: segments, Music: music, SFX: sfx, Active: true})

	}
	c.Status(200).JSON(fiber.Map{"profiles": profiles})
	return nil
}
func GetProfileById(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := int(claims["uid"].(float64))
	pid, err := strconv.Atoi(c.Params("id"))
	log.Debug().Int("UID", uid).Int("PID", pid).Send()
	GetRow, err := model.GetProfileById(uid, pid)
	if err != nil {
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}
	var oid float64
	var prid float64
	var name string
	var cursor string
	var bgType string
	var bgContent string
	var pcolour string
	var scolour string
	var segments float64
	var music bool
	var sfx bool
	var active bool
	prid, oid, name, cursor, bgType, bgContent, pcolour, scolour, segments, music, sfx, active, err = GetRow()
	// Exit if we get an error
	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}
	if !active {
		return c.SendStatus(204)
	}
	return c.JSON(fiber.Map{"profile": Profile{PID: prid, OID: oid, Name: name, Cursor: cursor, BGType: bgType, BGContent: bgContent, PColour: pcolour, SColour: scolour, Segments: segments, Music: music, SFX: sfx}})

}
func AddNewProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := int(claims["uid"].(float64))
	profile := new(Profile)
	if err := c.BodyParser(profile); err != nil {
		return err
	}
	m, err := model.GetProfileMaxID(uid)
	if err != nil {
		return err
	}
	err = model.AddProfile(m+1, uid, profile.Name, profile.Cursor, profile.BGType, profile.BGContent, profile.PColour, profile.SColour, profile.Segments, profile.Music, profile.SFX)
	if err != nil {
		return err
	}

	c.Status(201).JSON(fiber.Map{
		"pid": m + 1,
	})
	return nil
}
func UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := int(claims["uid"].(float64))
	pid, err := strconv.Atoi(c.Params("id"))
	profile := new(Profile)
	if err := c.BodyParser(profile); err != nil {
		return err
	}
	err = model.UpdateProfile(pid, uid, profile.Name, profile.Cursor, profile.BGType, profile.BGContent, profile.PColour, profile.SColour, profile.Segments, profile.Music, profile.SFX)
	if err != nil {
		return err
	}
	return c.SendStatus(202)
}
func DeleteProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := int(claims["uid"].(float64))
	pid, err := strconv.Atoi(c.Params("id"))
	err = model.DeleteProfile(pid, uid)
	if err != nil {
		return err
	}
	return c.SendStatus(202)
}
