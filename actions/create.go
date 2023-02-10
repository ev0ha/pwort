package actions

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ev0/pwort/encryption"
	"github.com/ev0/pwort/password"
	_ "github.com/mattn/go-sqlite3"
)

func Create(User *string, App *string, unsafe *bool) error {
	var userPwd string
	var appPwd string
	var userInput string

	scanner := bufio.NewScanner(os.Stdin)

	if *User != "" {
		if *App == "" {
			if *unsafe {
				fmt.Printf("Do you want to create new user %s with an unsafe password?\n", *User)
				fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
				fmt.Print("> ")
				scanner.Scan()
				userInput = scanner.Text()

				switch userInput {
				case "y":
					fmt.Println("Type the password:")
					fmt.Print("> ")
					scanner.Scan()
					userPwd = scanner.Text()
					createUnsafeUser(User, &userPwd)
				default:
					os.Exit(0)
				}
			} else {
				fmt.Printf("Do you want to create a new user %s?\n", *User)
				fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
				fmt.Print("> ")
				scanner.Scan()
				userInput = scanner.Text()

				switch userInput {
				case "y":
					fmt.Println("Type the password:")
					fmt.Print("> ")
					scanner.Scan()
					userPwd = scanner.Text()
					createUser(User, &userPwd)
				default:
					os.Exit(0)
				}
			}
		} else {
			if *unsafe {
				fmt.Printf("Do you want to create new app %s for user %s with an unsafe password?\n", *App, *User)
				fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
				fmt.Print("> ")
				scanner.Scan()
				userInput = scanner.Text()

				fmt.Println("Type the masterpassword:")
				fmt.Print("> ")
				scanner.Scan()
				userPwd = scanner.Text()

				if !password.CheckMasterPwd(User, &userPwd) {
					log.Fatal("Wrong masterpassword!")
				}

				fmt.Println("Type the password for the new app:")
				fmt.Print("> ")
				scanner.Scan()
				appPwd = scanner.Text()

				switch userInput {
				case "y":
					createUnsafeApp(User, App, &appPwd)
				default:
					os.Exit(0)
				}
			} else {
				fmt.Printf("Do you want to create new app %s for user %s?\n", *App, *User)
				fmt.Println("Type 'y' to continue. Otherwise type any other key to exit.")
				fmt.Print("> ")
				scanner.Scan()
				userInput = scanner.Text()

				fmt.Println("Type the masterpassword:")
				fmt.Print("> ")
				scanner.Scan()
				userPwd = scanner.Text()

				if !password.CheckMasterPwd(User, &userPwd) {
					log.Fatal("Wrong masterpassword!")
				}

				fmt.Println("Type the password for the new app:")
				fmt.Print("> ")
				scanner.Scan()
				appPwd = scanner.Text()

				if !password.CheckPwd(&appPwd) {
					log.Fatal("The specified password is not safe. Please use at least 12 characters, including lower/uppercase letters, numbers and special characters.")
				}

				switch userInput {
				case "y":
					createApp(User, App, &appPwd)
				default:
					os.Exit(0)
				}
			}
		}
	} else {
		fmt.Println("Please specify an user.")
	}
	return nil
}

func createUser(User *string, Pwd *string) {
	safe := password.CheckPwd(Pwd)

	if safe {
		db, err := sql.Open("sqlite3", "./passwords.db")
		if err != nil {
			log.Fatal("Could not connect to database")
		}
		defer db.Close()

		stmt, err := db.Prepare("INSERT OR IGNORE INTO users (name, pw) VALUES (?, ?)")
		if err != nil {
			db.Close()
			log.Fatal("Could not insert data into database.")
		}
		defer stmt.Close()

		hashedPwd, err := encryption.HashPwd(Pwd)
		if err != nil {
			log.Fatal("Could not create hash")
		}

		stmt.Exec(*User, hashedPwd)

		fmt.Printf("Created new user %s with password %s\n", *User, *Pwd)
	} else {
		log.Fatal("Please use at least one lowercase and uppercase letter, a digit and a special character [!§$%&/()=?\\+*~#'-_<>|,;.:]")
	}
}

func createUnsafeUser(User *string, Pwd *string) {
	db, err := sql.Open("sqlite3", "./passwords.db")
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT OR IGNORE INTO users (name, pw) VALUES (?, ?)")
	if err != nil {
		log.Fatal("Could not insert data into database.")
	}
	defer stmt.Close()

	hashedPwd, err := encryption.HashPwd(Pwd)
	if err != nil {
		log.Fatal("Could not create hash")
	}

	stmt.Exec(*User, hashedPwd)

	fmt.Printf("Created new user %s with unsafe password %s\n", *User, *Pwd)
}

func createApp(User *string, App *string, Pwd *string) {
	pwdCorrect := password.CheckMasterPwd(User, Pwd)

	if !pwdCorrect {
		log.Fatal("Wrong master password.")
	}

	safe := password.CheckPwd(Pwd)

	if safe {
		db, err := sql.Open("sqlite3", "./passwords.db")
		if err != nil {
			log.Fatal("Could not connect to database")
		}
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO apps (name, pw, user) VALUES (?, ?, ?)")
		if err != nil {
			log.Fatal("Could not insert data into database.")
		}
		defer stmt.Close()

		stmt.Exec(*App, *Pwd, *User)

		fmt.Printf("Created new app %s for user %s with password %s\n", *App, *User, *Pwd)
	} else {
		log.Fatal("Please use at least one lowercase and uppercase letter, a digit and a special character [!§$%&/()=?\\+*~#'-_<>|,;.:]")
	}
}

func createUnsafeApp(User *string, App *string, Pwd *string) {
	db, err := sql.Open("sqlite3", "./passwords.db")
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO apps (name, pw, user) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal("Could not insert data into database.")
	}
	defer stmt.Close()

	stmt.Exec(*App, *Pwd, *User)

	fmt.Printf("Created new app %s for user %s with unsafe password %s\n", *App, *User, *Pwd)
}
