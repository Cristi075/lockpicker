package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	expect "github.com/google/goexpect"
	"github.com/google/goterm/term"
)

const (
	timeout = 10 * time.Second

	checkCommand = "whoami"
	badResult    = "su: Authentication failure"
)

// Checks if a string is part of an array of strings
func ArrayContains(array []string, targetValue string) bool {
	for _, value := range array {
		if value == targetValue {
			return true
		}
	}
	return false
}

// Returns an array of all shells available on this system (from /etc/shells)
func GetShells() (result []string) {
	file, err := os.Open("/etc/shells")
	if err != nil {
		fmt.Println("ERROR. Could not read /etc/shells")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "/") {
			result = append(result, line)
		}
	}

	return result
}

// Returns an array of all usernames that exist on this system (from /etc/passwd)
func GetAllUsers() (result []string) {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		fmt.Println("ERROR. Could not read /etc/passwd")
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		result = append(result, line)
	}

	return
}

// Returns an array of usernames for users that have shells associated with them
func GetUsersWithShells(shells []string) (result []string) {
	users := GetAllUsers()
	for _, user := range users {
		fragments := strings.Split(user, ":")
		username := fragments[0]
		shell := fragments[len(fragments)-1]
		if ArrayContains(shells, shell) {
			//fmt.Println(username, "->", shell)
			result = append(result, username)
		}
	}
	return
}

// Tries to login with the given credentials (username/password)
// Returns the result through resultChannel
func CheckUserPassword(username string, password string, resultChannel chan Result) {
	tmp, _, err := expect.Spawn(fmt.Sprintf("su %s -c %s", username, checkCommand), -1)
	if err != nil {
		fmt.Println("ERROR. Could not spawn terminal")
		//os.Exit(3)
	}
	defer tmp.Close()

	passwordRe := regexp.MustCompile("Password:")
	promptRe := regexp.MustCompile("%")

	tmp.Expect(passwordRe, timeout)
	tmp.Send(password + "\n")

	result, _, _ := tmp.Expect(promptRe, timeout)
	result = strings.TrimSpace(result)

	if result == badResult {
		resultChannel <- Result{username, false, ""}
	} else {
		if result == username {
			resultChannel <- Result{username, true, password}
		}
	}

	//fmt.Println("Error while checking user", username)
	//fmt.Println("Received invalid output for 'whoami' command:")
	//fmt.Println(result)
	resultChannel <- Result{username, false, ""}
}

// Checks if the given username has simple passwords
// The password candidates are obtained by calling the GeneratePasswords function
func CheckUser(username string, resultChannel chan Result) {
	passwordList := GeneratePasswords(username)
	nrOfPasswords := len(passwordList)

	resultChannel2 := make(chan Result, 10)

	for _, password := range passwordList {
		go CheckUserPassword(username, password, resultChannel2)
	}

	for i := 0; i < nrOfPasswords; i++ {
		res := <-resultChannel2

		if res.found {
			resultChannel <- Result{username, true, res.password}
			return
		}
	}

	resultChannel <- Result{username, false, ""}
}

type Result struct {
	username string
	found    bool
	password string
}

// Generates password candidates for the given username
// TODO: Let the user chose which rules are used for this
// TODO: Implement more rules
func GeneratePasswords(username string) (result []string) {
	// Password = identical to username. Always on
	result = append(result, username)

	// Password = username + 1 digit. On by default
	for i := 0; i < 10; i++ {
		result = append(result, fmt.Sprintf("%s%d", username, i))
	}

	// Password = username + 2 digit. On by default
	for i := 0; i < 100; i++ {
		result = append(result, fmt.Sprintf("%s%02d", username, i))
	}
	// Password = username + 3 digit. Off by default
	/*
		for i := 0; i < 999; i++ {
			result = append(result, fmt.Sprintf("%s%03d", username, i))
		}
	*/

	return
}

func PrintBanner() {
	fmt.Println()
	fmt.Println(term.BBlue("---------------------------------------------------------"))
	fmt.Println(`
     _            _          _      _             
    | |          | |        (_)    | |            
    | | ___   ___| | ___ __  _  ___| | _____ _ __ 
    | |/ _ \ / __| |/ / '_ \| |/ __| |/ / _ \ '__|
    | | (_) | (__|   <| |_) | | (__|   <  __/ |   
    |_|\___/ \___|_|\_\ .__/|_|\___|_|\_\___|_|   
                      | |                         
                      |_|                 
                    > github.com/Cristi075/lockpicker`)
	fmt.Println(term.BBlue("---------------------------------------------------------"))
}

func main() {
	PrintBanner()

	shells := GetShells()
	fmt.Println("Found", len(shells), "shells")

	users := GetUsersWithShells(shells)
	nrOfUsers := len(users)
	fmt.Println("Found", nrOfUsers, "users with shells:", users)

	results := make(chan Result, nrOfUsers)

	for _, username := range users {
		go CheckUser(username, results)
	}

	for i := 0; i < nrOfUsers; i++ {
		res := <-results

		if res.found {
			fmt.Printf("Checking user %s", res.username)
			fmt.Print(term.Greenf(" - Valid password found." + " Password: " + res.password + "\n"))
		} else {
			fmt.Printf("Checking user %s", res.username)
			fmt.Print(term.Redf(" - No password found\n"))
		}
	}
}
