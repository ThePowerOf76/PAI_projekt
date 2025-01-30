package model

func GetProfiles(uid int) (func() (float64, float64, string, string, string, string, string, string, float64, bool, bool, bool, error), error) {
	rows, err := DB.Query("SELECT * FROM profiles WHERE ownerid = ?", uid)

	return func() (float64, float64, string, string, string, string, string, string, float64, bool, bool, bool, error) {
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
		if rows.Next() {
			err := rows.Scan(&pid, &oid, &name, &cursor, &bgType, &bgContent, &pcolour, &scolour, &segments, &music, &sfx, &active)
			return pid, oid, name, cursor, bgType, bgContent, pcolour, scolour, segments, music, sfx, active, err
		} else {
			rows.Close()
			return -1, -1, "", "", "", "", "", "", -1, false, false, false, nil
		}

	}, err
}
func GetProfileById(uid int, pid int) (func() (float64, float64, string, string, string, string, string, string, float64, bool, bool, bool, error), error) {
	rows, err := DB.Query("SELECT * FROM profiles WHERE pid = ? AND ownerid = ?", pid, uid)

	return func() (float64, float64, string, string, string, string, string, string, float64, bool, bool, bool, error) {
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
		if rows.Next() {
			err := rows.Scan(&pid, &oid, &name, &cursor, &bgType, &bgContent, &pcolour, &scolour, &segments, &music, &sfx, &active)
			return pid, oid, name, cursor, bgType, bgContent, pcolour, scolour, segments, music, sfx, active, err
		} else {
			rows.Close()
			return -1, -1, "", "", "", "", "", "", -1, false, false, false, nil
		}

	}, err
}
func GetProfileMaxID(ownerid int) (int, error) {
	rows, err := DB.Query("SELECT MAX(pid) FROM profiles WHERE ownerid = ?", ownerid)
	if err != nil {
		return -1, err
	}
	var m int
	rows.Next()
	rows.Scan(&m)
	return m, nil
}
func AddProfile(pid int, uid int, Name string, Cursor string, BGType string, BGContent string, PColour string, SColour string, Segments float64, Music bool, SFX bool) error {
	_, err := DB.Exec("INSERT INTO profiles VALUES(?,?,?,?,?,?,?,?,?,?,?,?)", pid, uid, Name, Cursor, BGType, BGContent, PColour, SColour, Segments, Music, SFX, true)
	return err
}
func UpdateProfile(pid int, uid int, Name string, Cursor string, BGType string, BGContent string, PColour string, SColour string, Segments float64, Music bool, SFX bool) error {
	_, err := DB.Exec("UPDATE profiles SET name=?, crsr=?, bgtype=?, bgcontent=?, pcolour=?, scolour=?, segments=?, music=?, sfx=? WHERE pid = ? AND ownerid = ?", Name, Cursor, BGType, BGContent, PColour, SColour, Segments, Music, SFX, pid, uid)
	return err
}
func DeleteProfile(pid int, uid int) error {
	_, err := DB.Exec("UPDATE profiles SET active=false WHERE pid = ? AND ownerid = ?", pid, uid)
	return err
}
