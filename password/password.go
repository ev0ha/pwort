package password

import (
	"database/sql"
	"log"
	"strings"

	"github.com/ev0/pwort/encryption"
)

func CheckPwd(Pwd *string) bool {
	if len(*Pwd) < 12 {
		log.Fatal("Password too short. Please use a password with at least 12 characters.")
	} else {
		lowerCase := strings.ContainsAny(*Pwd, "abcdefghijklmnopqrstuvwxyz")
		upperCase := strings.ContainsAny(*Pwd, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		digits := strings.ContainsAny(*Pwd, "0123456789")
		special := strings.ContainsAny(*Pwd, "!§$%&/()=?\\+*~#'-_<>|,;.:")
		if lowerCase && upperCase && digits && special {
			return true
		} else {
			return false
		}
	}
	return true
}

func CheckMasterPwd(user *string, pwd *string) bool {
	var hashedPwd string

	db, err := sql.Open("sqlite3", "./passwords.db")
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	defer db.Close()

	sqlStmt := `SELECT pw FROM users WHERE name = ?`

	row := db.QueryRow(sqlStmt, *user)

	switch err := row.Scan(&hashedPwd); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return encryption.CheckPwdHash(pwd, &hashedPwd)
	default:
		panic(err)
	}
}
