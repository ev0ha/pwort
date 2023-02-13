package actions

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ev0/pwort/password"
)

func Delete(User *string, App *string) error {
	var userInput string
	scanner := bufio.NewScanner(os.Stdin)

	if *User != "" {
		if *App == "" {
			fmt.Printf("Do you want to delete user %s and all apps belonging to that user?\n", *User)
			fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
			fmt.Print("> ")
			scanner.Scan()
			userInput = scanner.Text()

			switch userInput {
			case "y":
				deleteUser(User)
			default:
				os.Exit(0)
			}
		} else {
			fmt.Printf("Do you want to delete app %s for user %s?\n", *App, *User)
			fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
			fmt.Print("> ")
			scanner.Scan()
			userInput = scanner.Text()

			switch userInput {
			case "y":
				deleteApp(User, App)
			default:
				os.Exit(0)
			}
		}
	} else {
		fmt.Println("Please specify an user and optionally an app.")
	}
	return nil
}

func deleteUser(User *string) {
	scanner := bufio.NewScanner(os.Stdin)
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

	_, err = db.Exec("DELETE FROM users WHERE name = ?", *User)
	if err != nil {
		log.Fatal("Error deleting user.")
	}

	_, err = db.Exec("DELETE FROM apps WHERE user = ?", *User)
	if err != nil {
		log.Fatal("Error deleting apps of specified user.")
	}

	fmt.Printf("Deleted user %s and all their apps.\n", *User)
}

func deleteApp(User *string, App *string) {
	scanner := bufio.NewScanner(os.Stdin)
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

	_, err = db.Exec("DELETE FROM apps WHERE name = ? AND user = ?", *App, *User)
	if err != nil {
		log.Fatal("Error deleting app for specified user.")
	}

	fmt.Printf("Deleted app %s for user %s\n", *App, *User)
}
