package model

func GetScores() (func() (string, string, string, float64, error), error) {
	rows, err := DB.Query("SELECT crsr, pcolour, nickname, score FROM scores sc LEFT JOIN profiles pr ON sc.profile = pr.pid AND sc.`user` = pr.ownerid LEFT JOIN users us ON pr.ownerid = us.uid ORDER BY score DESC")

	return func() (string, string, string, float64, error) {
		var name string
		var cursor string
		var color string
		var score float64
		if rows.Next() {
			err := rows.Scan(&cursor, &color, &name, &score)
			return cursor, color, name, score, err
		} else {
			rows.Close()
			return "", "", "", -1, nil
		}

	}, err

}
func AddNewScore(player_id float64, profile_id float64, score float64) error {
	_, err := DB.Exec("INSERT INTO scores(score_id, user, profile, score) VALUES(NULL,?,?,?)", player_id, profile_id, score)
	return err
}

/*
import (

	"database/sql"
	b64 "encoding/base64"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"golang.org/x/crypto/argon2"
	"os"
	"time"

)
*/
