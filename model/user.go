package model

import "fmt"

func GetUserCredentials(username string) (int, string, error) {
	rows, err := DB.Query("SELECT uid, password FROM users WHERE nickname = ?", username)
	defer rows.Close()
	if err != nil {
		return 0, "", err
	}
	var userid int
	var password string
	rows.Next()
	rows.Scan(&userid, &password)
	return userid, password, nil
}
func AddNewToken(token string, expiry string) error {
	_, err := DB.Exec("INSERT INTO tokens VALUES(NULL, ?, ?)", token, expiry)
	return err
}
func InvalidateToken(token string) error {
	_, err := DB.Exec("DELETE FROM tokens WHERE token = ?", token)
	return err

}
func GetUserCountByName(username string) (int, error) {
	row, err := DB.Query(fmt.Sprintf("SELECT COUNT(*) FROM users WHERE nickname = '%s'", username))
	defer row.Close()
	var r int
	if err != nil {
		return -1, err
	}

	row.Next()
	err = row.Scan(&r)
	if err != nil {
		return -1, err
	}
	return r, nil
}
func AddNewUser(username string, passwordhash string) error {
	_, err := DB.Exec("INSERT INTO users VALUES(NULL, ?, ?)", username, passwordhash)
	return err
}
