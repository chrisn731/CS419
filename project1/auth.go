package main

import (
	"encoding/json"
	"os"
	"log"
	"io/ioutil"
	"fmt"
)

const (
	usersFile = "./users.json"
	domainsFile = "./domains.json"
	typesFile = "./types.json"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Domain struct {
	Users []string
	Operations []string
}

var domains2 = make(map[string]Domain)
var domains = make(map[string][]string)
var users []User = make([]User, 0)
var types = make(map[string][]string)

func printUsage() {
	fmt.Print("Usage: ./auth <command> [options...]\n",
		"Commands:\n",
		"  AddUser\n",
		"  etc\n",
	)
	os.Exit(0)
}

func sync() {
	syncContainer := func(file string, container interface{}) {
		bytes, err := json.MarshalIndent(container, "", " ")
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(file, bytes, 0644)
	}
	syncContainer(usersFile, users)
	syncContainer(domainsFile, domains)
	syncContainer(typesFile, types)
}

func fetchTables() {
	fetch := func(file string, container interface{}) {
		if _, err := os.Stat(file); err != nil {
			return
		}
		fbytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(fbytes, container)
		if err != nil {
			log.Fatal(err)
		}
	}
	fetch(usersFile, &users)
	fetch(domainsFile, &domains)
	fetch(typesFile, &types)
}

// Defines a new user with a password.
// The username CANNOT be an empty string. The password CAN be an empty string.
func addUser(username, pass string) string {
	if username == "" {
		return "Error: username missing"
	}
	for _, elm := range users {
		if elm.Username == username {
			return "Error: user exists"
		}
	}
	users = append(users, User{
		Username: username,
		Password: pass,
	})
	sync()
	return "Success"
}

// Validates a user's password by username and password.
func authenticate(username, pass string) string {
	for _, elm := range users {
		if elm.Username == username {
			if elm.Password == pass {
				return "Success"
			}
			return "Error: bad password"
		}
	}
	return "Error: no such user"
}

// Assign a user to a domain.
// If the domain name does not exist, it will be created.
// If a user does not exist, this function will return an error.
// A user may belong to multiple domains, but domains can not have duplicates.
// The domain name must be a non-empty string.
func setDomain(user, dName string) string {
	v, ok := domains[dName]
	if !ok {
		v = make([]string, 0)
		domains[dName] = v
	}
	for _, elm := range v {
		if elm == user {
			return "Error: user exists"
		}
	}
	for _, elm := range users {
		if elm.Username == user {
			v = append(v, user)
			domains[dName] = v
			return "Success"

		}
	}
	return "Error: no such user"
}

// List all the users within a domain.
// User output will be newline seperated.
// If the domain does not exist OR the domain has no users, there will be no output.
// The passed in domain name must NOT be empty.
func domainInfo(args []string) {
	if len(args) > 1 {
		fmt.Println("Error: too many arguments for DomainInfo")
		return
	} else if len(args) < 1 {
		fmt.Println("Error: too few arguments for DomainInfo")
		return
	}

	dName := args[0]
	if dName == "" {
		fmt.Println("Error missing name of domain")
	}

	if v, ok := domains[dName]; ok {
		for _, name := range v {
			fmt.Println(name)
		}
	}
}

func setType(args []string) {
	if len(args) < 2 {
		fmt.Println("Error: not enough arguments to SetType!")
		return
	} else if len(args) > 2 {
		fmt.Println("Error: too many arguments to SetType!")
		return
	}

	objName, typeName := args[0], args[1]
	if objName == "" {
		fmt.Println("Error: objectname is empty")
		return
	} else if typeName == "" {
		fmt.Println("Error: type is empty")
		return
	}

	objs, ok := types[typeName]
	if !ok {
		objs = make([]string, 0)
	}
	objs = append(objs, objName)
	types[typeName] = objs
	fmt.Println("Success")
}

func typeInfo(args []string) {
	if len(args) > 1 {
		fmt.Println("Error: too many arguments for TypeInfo")
		return
	} else if len(args) < 1 {
		fmt.Println("Error: not enough arguments for TypeInfo")
		return
	}

	typeName := args[0]
	if typeName == "" {
		fmt.Println("Error: type is empty")
		return
	}

	objs, ok := types[typeName]
	if !ok || len(objs) == 0 {
		return
	}

	for _, obj := range objs {
		fmt.Println(obj)
	}
}

func addAccess(args []string) {

}

func canAccess(args []string) {

}

func cleanup_and_exit() {
	sync()
	os.Exit(0)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
	}

	fetchTables()
	switch command, cargs := args[0], args[1:]; command {
	case "AddUser":
		if len(cargs) != 2 {
			printUsage()
		}
		ret := addUser(args[1], args[2])
		fmt.Println(ret)
	case "Authenticate":
		if len(cargs) != 2 {
			printUsage()
		}
		ret := authenticate(args[1], args[2])
		fmt.Println(ret)
	case "SetDomain":
		if len(cargs) != 2 {
			printUsage()
		}
		ret := setDomain(args[1], args[2])
		fmt.Println(ret)
	case "DomainInfo":
		if len(cargs) != 1 {
			printUsage()
		}
		if args[1] == "" {
			return
		}
		domainInfo(cargs)
	case "SetType":
		setType(cargs)
	case "TypeInfo":
		typeInfo(cargs)
	case "AddAccess":
		addAccess(cargs)
	case "CanAccess":
		canAccess(cargs)
	default:
		printUsage()
	}
	cleanup_and_exit()
}
