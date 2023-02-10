package actions

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ev0/pwort/password"
)

func Show(User *string, App *string) error {
	if *User != "" {
		if *App == "" {
			showUser(User)
		} else {
			showApp(User, App)
		}
	} else {
		fmt.Println("Please specify an user and optionally an app.")
	}
	return nil
}

func showUser(User *string) {
	//var name string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please type the masterpassword for the specified user.")
	fmt.Print("> ")
	scanner.Scan()
	mpwd := scanner.Text()

	pwdCorrect := password.CheckMasterPwd(User, &mpwd)
	if !pwdCorrect {
		log.Fatal("Wrong master password or user doesn't exist.")
	} else {
		fmt.Printf("User %s exists.\n", *User)
	}
}

func showApp(User *string, App *string) {
	var pwd string
	var userInput string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Do you want the password of app %s for user %s to be displayed?\nPress 'y' for yes or anything else if not.\n", *App, *User)
	fmt.Print("> ")
	scanner.Scan()
	userInput = scanner.Text()
	if userInput != "y" {
		log.Fatal("Password will not be displayed.")
	}

	fmt.Println("Please type the masterpassword for the specified user.")
	fmt.Print("> ")
	scanner.Scan()
	mpwd := scanner.Text()

	pwdCorrect := password.CheckMasterPwd(User, &mpwd)
	if !pwdCorrect {
		log.Fatal("Wrong master password.")
	}

	db, err := sql.Open("sqlite3", "./passwords.db")
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	defer db.Close()

	sqlStmt := `SELECT pw FROM apps WHERE name = ? AND user = ?`

	row := db.QueryRow(sqlStmt, *App, *User)

	switch err := row.Scan(&pwd); err {
	case sql.ErrNoRows:
		fmt.Println("Specified app and/or user doesn't exist!")
	case nil:
		fmt.Printf("The password for app %s for user %s is: %s\n", *App, *User, pwd)
	default:
		panic(err)
	}

}
