package actions

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ev0/pwort/encryption"
	"github.com/ev0/pwort/password"
)

func Update(User *string, App *string, unsafe *bool) error {
	var userInput string
	var mpwd string
	scanner := bufio.NewScanner(os.Stdin)

	if *User != "" {
		if *App == "" {
			if *unsafe {
				fmt.Printf("Do you want to update user %s with an unsafe password?\n", *User)
				fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
				fmt.Print("> ")
				scanner.Scan()
				userInput = scanner.Text()

				fmt.Println("Please type the masterpassword:")
				fmt.Print("> ")
				scanner.Scan()
				mpwd = scanner.Text()
				if !password.CheckMasterPwd(User, &mpwd) {
					log.Fatal("Wrong masterpassword!")
				}

				switch userInput {
				case "y":
					updateUnsafeUser(User)
				default:
					os.Exit(0)
				}
			} else {
				fmt.Printf("Do you want to update user %s with a new password?\n", *User)
				fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
				fmt.Print("> ")
				scanner.Scan()
				userInput = scanner.Text()

				fmt.Println("Please type the masterpassword:")
				fmt.Print("> ")
				scanner.Scan()
				mpwd = scanner.Text()
				if !password.CheckMasterPwd(User, &mpwd) {
					log.Fatal("Wrong masterpassword!")
				}

				switch userInput {
				case "y":
					updateUser(User)
				default:
					os.Exit(0)
				}
			}
		} else {
			if *unsafe {
				fmt.Printf("Do you want to update app %s for user %s with an unsafe password?\n", *App, *User)
				fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
				fmt.Print("> ")
				scanner.Scan()
				userInput = scanner.Text()

				fmt.Println("Please type the masterpassword:")
				fmt.Print("> ")
				scanner.Scan()
				mpwd = scanner.Text()
				if !password.CheckMasterPwd(User, &mpwd) {
					log.Fatal("Wrong masterpassword!")
				}

				switch userInput {
				case "y":
					updateUnsafeApp(User, App)
				default:
					os.Exit(0)
				}
			} else {
				fmt.Printf("Do you want to update the app %s for user %s?\n", *App, *User)
				fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
				fmt.Print("> ")
				scanner.Scan()
				userInput = scanner.Text()

				fmt.Println("Please type the masterpassword:")
				fmt.Print("> ")
				scanner.Scan()
				mpwd = scanner.Text()
				if !password.CheckMasterPwd(User, &mpwd) {
					log.Fatal("Wrong masterpassword!")
				}

				switch userInput {
				case "y":
					updateApp(User, App)
				default:
					os.Exit(0)
				}
			}
		}
	} else {
		fmt.Println("Please specify an user and optionally an app.")
	}
	return nil
}

func updateUser(User *string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter new username or press enter if you don't want to change the username:")
	fmt.Print("> ")
	scanner.Scan()
	newName := scanner.Text()
	fmt.Println("Enter new password or press enter if you don't want to change the password:")
	fmt.Print("> ")
	scanner.Scan()
	newPassword := scanner.Text()

	if len(newPassword) > 0 {
		if password.CheckPwd(&newPassword) {
			if len(newName) > 0 {
				// update username and password
				db, err := sql.Open("sqlite3", "./passwords.db")
				if err != nil {
					log.Fatal("Could not connect to database")
				}
				defer db.Close()

				_, err = db.Exec("UPDATE users SET name = ?, pw = ? WHERE name = ?", newName, newPassword, User)
				if err != nil {
					log.Fatal("Error updating.")
				}

				fmt.Printf("Updated user %s with new name %s and new password %s.\n", *User, newName, newPassword)
			} else {
				// update only password
				db, err := sql.Open("sqlite3", "./passwords.db")
				if err != nil {
					log.Fatal("Could not connect to database")
				}
				defer db.Close()

				_, err = db.Exec("UPDATE users SET pw = ? WHERE name = ?", newPassword, User)
				if err != nil {
					log.Fatal("Error updating.")
				}

				fmt.Printf("Updated user %s with new password %s.\n", *User, newPassword)
			}
		} else {
			log.Fatal("The password is not safe. Please use at least 12 characters, including at least one lower/uppercase letters, numbers and special characters")
		}
	} else {
		if len(newName) > 0 {
			// update username
			db, err := sql.Open("sqlite3", "./passwords.db")
			if err != nil {
				log.Fatal("Could not connect to database")
			}
			defer db.Close()

			_, err = db.Exec("UPDATE users SET name = ? WHERE name = ?", newName, User)
			if err != nil {
				log.Fatal("Error updating.")
			}

			fmt.Printf("Updated user %s with new name %s.\n", *User, newName)
		} else {
			log.Fatal("You neither specified a new name nor a new password!")
		}
	}
}

func updateUnsafeUser(User *string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter new username or press enter if you don't want to change the username:")
	fmt.Print("> ")
	scanner.Scan()
	newName := scanner.Text()
	fmt.Println("Enter new UNSAFE password or press enter if you don't want to change the password:")
	fmt.Print("> ")
	scanner.Scan()
	newPassword := scanner.Text()

	if newPassword != "" {
		if newName != "" {
			// update username and password
			db, err := sql.Open("sqlite3", "./passwords.db")
			if err != nil {
				log.Fatal("Could not connect to database")
			}
			defer db.Close()

			hashedPwd, err := encryption.HashPwd(&newPassword)
			if err != nil {
				log.Fatal("Could not create hash.")
			}
			_, err = db.Exec("UPDATE users SET name = ?, pw = ? WHERE name = ?", hashedPwd, *User)
			if err != nil {
				log.Fatal("Could not create hash.")
			}

			fmt.Printf("Updated user %s with new name %s and new password %s.\n", *User, newName, newPassword)
		} else {
			// update only password
			db, err := sql.Open("sqlite3", "./passwords.db")
			if err != nil {
				log.Fatal("Could not connect to database")
			}
			defer db.Close()

			hashedPwd, err := encryption.HashPwd(&newPassword)
			if err != nil {
				log.Fatal("Could not create hash.")
			}

			db.Exec("UPDATE users SET pw = ? WHERE name = ?", hashedPwd, *User)

			fmt.Printf("Updated user %s with new password %s.\n", *User, newPassword)
		}
	} else {
		if newName != "" {
			// update only username
			db, err := sql.Open("sqlite3", "./passwords.db")
			if err != nil {
				log.Fatal("Could not connect to database")
			}
			defer db.Close()

			_, err = db.Exec("UPDATE users SET name = ? WHERE name = ?", newName, *User)
			if err != nil {
				log.Fatal("Error updating.")
			}

			fmt.Printf("Updated user %s with new name %s.\n", *User, newName)
		} else {
			log.Fatal("You neither specified a new name nor a new password!")
		}
	}
}

func updateApp(User *string, App *string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter new appname or press enter if you don't want to change the name of the app:")
	fmt.Print("> ")
	scanner.Scan()
	newName := scanner.Text()
	fmt.Println("Enter new password or press enter if you don't want to change the password:")
	fmt.Print("> ")
	scanner.Scan()
	newPassword := scanner.Text()

	if newPassword != "" {
		if password.CheckPwd(&newPassword) {
			if newName != "" {
				// update name of app app and password
				db, err := sql.Open("sqlite3", "./passwords.db")
				if err != nil {
					log.Fatal("Could not connect to database")
				}
				defer db.Close()

				_, err = db.Exec("UPDATE apps SET name = ?, pw = ? WHERE name = ? AND user = ?", newName, newPassword, *App, *User)
				if err != nil {
					log.Fatal("Error updating.")
				}

				fmt.Printf("Updated app %s with new name %s and new password %s.\n", *App, newName, newPassword)
			} else {
				// update only password
				db, err := sql.Open("sqlite3", "./passwords.db")
				if err != nil {
					log.Fatal("Could not connect to database")
				}
				defer db.Close()

				_, err = db.Exec("UPDATE apps SET pw = ? WHERE name = ? AND user = ?", newPassword, *App, *User)
				if err != nil {
					log.Fatal("Error updating.")
				}

				fmt.Printf("Updated app %s with new password %s.\n", *App, newPassword)
			}
		} else {
			log.Fatal("The password is not safe. Please use at least 12 characters, including at least one lower/uppercase letters, numbers and special characters")
		}
	} else {
		if newName != "" {
			// update name of app
			db, err := sql.Open("sqlite3", "./passwords.db")
			if err != nil {
				log.Fatal("Could not connect to database")
			}
			defer db.Close()

			_, err = db.Exec("UPDATE apps SET name = ? WHERE name = ? AND user = ?", newName, *App, *User)
			if err != nil {
				log.Fatal("Error updating.")
			}

			fmt.Printf("Updated app %s with new name %s.\n", *App, newName)
		} else {
			log.Fatal("You neither specified a new name nor a new password!")
		}
	}
}

func updateUnsafeApp(User *string, App *string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter new appname or press enter if you don't want to change the name of the app:")
	fmt.Print("> ")
	scanner.Scan()
	newName := scanner.Text()
	fmt.Println("Enter new UNSAFE password or press enter if you don't want to change the password:")
	fmt.Print("> ")
	scanner.Scan()
	newPassword := scanner.Text()

	if newPassword != "" {
		if newName != "" {
			// update name of app and password
			db, err := sql.Open("sqlite3", "./passwords.db")
			if err != nil {
				log.Fatal("Could not connect to database")
			}
			defer db.Close()

			_, err = db.Exec("UPDATE apps SET name = ?, pw = ? WHERE name = ? AND user = ?", newName, newPassword, *App, *User)
			if err != nil {
				log.Fatal("Error updating.")
			}

			fmt.Printf("Updated app %s with new name %s and new password %s.\n", *App, newName, newPassword)
		} else {
			// update only password
			db, err := sql.Open("sqlite3", "./passwords.db")
			if err != nil {
				log.Fatal("Could not connect to database")
			}
			defer db.Close()

			_, err = db.Exec("UPDATE apps SET pw = ? WHERE name = ? AND user = ?", newPassword, *App, *User)
			if err != nil {
				log.Fatal("Error updating.")
			}

			fmt.Printf("Updated app %s with new password %s.\n", *App, newPassword)
		}
	} else {
		if newName != "" {
			// update name of app
			db, err := sql.Open("sqlite3", "./passwords.db")
			if err != nil {
				log.Fatal("Could not connect to database")
			}
			defer db.Close()

			_, err = db.Exec("UPDATE apps SET name = ? WHERE name = ? AND user = ?", newName, *App, *User)
			if err != nil {
				log.Fatal("Error updating.")
			}

			fmt.Printf("Updated app %s with new name %s.\n", *App, newName)
		} else {
			log.Fatal("You neither specified a new name nor a new password!")
		}
	}
}
