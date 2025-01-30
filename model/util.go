package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"os"
)

var DB *sql.DB

func DBInit() error {
	var err error
	dsn := "adam:haslomaslo@/gdoc"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	DB = sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter)

	if err = DB.Ping(); err != nil {
		return err
	}
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS users (uid INTEGER PRIMARY KEY AUTO_INCREMENT, nickname VARCHAR(100) UNIQUE NOT NULL, password VARCHAR(255) NOT NULL) ")
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS tokens (id INTEGER PRIMARY KEY AUTO_INCREMENT, token VARCHAR(300) NOT NULL, expires DATETIME NOT NULL)")
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS profiles (pid INTEGER NOT NULL, ownerid INTEGER NOT NULL,
			 name VARCHAR(20) NOT NULL, crsr CHAR(8) NOT NULL, bgType CHAR(10) NOT NULL,
			 bgcontent VARCHAR(50) NOT NULL, pcolour CHAR(7) NOT NULL, scolour CHAR(7) NOT NULL,
			 segments INTEGER DEFAULT 30, music BOOLEAN DEFAULT false, sfx BOOLEAN DEFAULT false, active BOOLEAN DEFAULT true,
		  FOREIGN KEY (ownerid) REFERENCES users(uid) ON DELETE CASCADE, CONSTRAINT ProfileUserPK PRIMARY KEY (pid, ownerid))`)
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scores (score_id INTEGER PRIMARY KEY AUTO_INCREMENT,
			 profile INTEGER NOT NULL, user INTEGER NOT NULL, score INTEGER NOT NULL,
		  FOREIGN KEY (profile) REFERENCES profiles(pid) ON DELETE CASCADE, FOREIGN KEY (user) REFERENCES profiles(ownerid))`)
	if err != nil {
		panic(err)
	}
	return nil

}
func CheckTokenValidity(token string) (int, error) {
	row := DB.QueryRow("SELECT COUNT(*) from tokens where token = ?", token)
	var len int
	err := row.Scan(&len)
	if err != nil {
		return -1, err
	}
	return len, nil
}
