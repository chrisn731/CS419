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

var domains = make(map[string][]string)
var users []User = make([]User, 0)

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
}

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

func setDomain(user, dName string) string {
	v, ok := domains[dName]
	if !ok {
		v = make([]string, 0)
		domains[dName] = v
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


func domainInfo(dName string) string {
	if dName == "" {
		return "Error missing domain"
	}

	if v, ok := domains[dName]; ok {
		for _, name := range v {
			fmt.Println(name)
		}
	}
	return ""
}

func cleanup_and_exit() {
	sync()
	os.Exit(0)
}

func main() {
	args := os.Args[1:]

	fetchTables()
	switch command := args[0]; command {
	case "AddUser":
		if len(args[1:]) != 2 {
			printUsage()
		}
		ret := addUser(args[1], args[2])
		fmt.Println(ret)
		cleanup_and_exit()
	case "Authenticate":
		if len(args[1:]) != 2 {
			printUsage()
		}
		ret := authenticate(args[1], args[2])
		fmt.Println(ret)
		cleanup_and_exit()
	case "SetDomain":
		if len(args[1:]) != 2 {
			printUsage()
		}
		ret := setDomain(args[1], args[2])
		fmt.Println(ret)
		cleanup_and_exit()
	case "DomainInfo":
		if len(args[1:]) != 1 {
			printUsage()
		}
		if args[1] == "" {
			return
		}
		domainInfo(args[1])
		cleanup_and_exit()
	case "SetType":
	case "TypeInfo":
	case "AddAccess":
	case "CanAccess":
	default:
		printUsage()
	}
	printUsage()
}
