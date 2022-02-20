package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	usersFile = "./users.json"
	domainsFile = "./domains.json"
	typesFile = "./types.json"
)

// Represents the User in our authentication system.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domains []string
}

// Represents a permission granted to a domain to a type
type Permission struct {
	// The operation allowed on the type
	Operation string

	// The type allowed to be accessed.
	Type string
}

// Represents the domains in our authentication system
type Domain struct {
	// A named reference to the users within the domain
	Users []string

	// Permissions that belong to this domain
	Permissions []Permission
}

var (
	domains = make(map[string]Domain)
	users []User = make([]User, 0)
	types = make(map[string][]string)
)

func die(errm string, args ...interface{}) {
	fmt.Printf(errm, args...)
	fmt.Println("")
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
func addUser(args []string) {
	if len(args) > 2 {
		die("Error: too many arguments for AddUser")
	} else if len(args) < 2 {
		die("Error: missing operands for AddUser")
	}

	username, pass := args[0], args[1]
	if username == "" {
		die("Error: username missing")
	}
	for _, elm := range users {
		if elm.Username == username {
			die("Error: user exists")
		}
	}
	users = append(users, User{
		Username: username,
		Password: pass,
	})
	fmt.Println("Success")
}

// Validates a user's password by username and password.
func authenticate(args []string) {
	if len(args) > 2 {
		die("Error: too many arguments for Authenticate")
	} else if len(args) < 2 {
		die("Error: missing operands for Authenticate")
	}

	username, pass := args[0], args[1]
	for _, elm := range users {
		if elm.Username == username {
			if elm.Password == pass {
				fmt.Println("Success")
				return
			}
			die("Error: bad password")
		}
	}
	die("Error: no such user")
}

// Assign a user to a domain.
// If the domain name does not exist, it will be created.
// If a user does not exist, this function will return an error.
// A user may belong to multiple domains, but domains can not have duplicates.
// The domain name must be a non-empty string.
func setDomain(args []string) {
	if len(args) > 2 {
		die("Error: too many arguments for SetDomain")
	} else if len(args) < 2 {
		die("Error: missing operands for SetDomain")
	}

	user, dName := args[0], args[1]
	if dName == "" {
		die("Error: missing domain")
	}

	v, ok := domains[dName]
	if !ok {
		v = Domain{}
		domains[dName] = v
	}
	// Go over the domains users
	for _, elm := range v.Users {
		if elm == user {
			die("Error: user exists")
		}
	}
	// Update the user
	for idx, elm := range users {
		if elm.Username == user {
			v.Users = append(v.Users, user)
			users[idx].Domains = append(users[idx].Domains, dName)
			domains[dName] = v
			fmt.Println("Success")
			return

		}
	}
	die("Error: no such user")
}

// List all the users within a domain.
// User output will be newline seperated.
// If the domain does not exist OR the domain has no users, there will be no output.
// The passed in domain name must NOT be empty.
func domainInfo(args []string) {
	if len(args) > 1 {
		die("Error: too many arguments for DomainInfo")
	} else if len(args) < 1 {
		die("Error: missing operands for DomainInfo")
	}

	dName := args[0]
	if dName == "" {
		die("Error: missing domain")
	}

	if v, ok := domains[dName]; ok {
		for _, name := range v.Users {
			fmt.Println(name)
		}
	}
}

func setType(args []string) {
	if len(args) < 2 {
		die("Error: not enough arguments to SetType!")
	} else if len(args) > 2 {
		die("Error: too many arguments to SetType!")
	}

	objName, typeName := args[0], args[1]
	if objName == "" {
		die("Error: objectname is empty")
	} else if typeName == "" {
		die("Error: type is empty")
	}

	objs, ok := types[typeName]
	if !ok {
		objs = make([]string, 0)
	}
	newObj := true
	for _, obj := range objs {
		if obj == objName {
			newObj = false
			break
		}
	}
	if newObj {
		objs = append(objs, objName)
		types[typeName] = objs
	}
	fmt.Println("Success")
}

func typeInfo(args []string) {
	if len(args) > 1 {
		die("Error: too many arguments for TypeInfo")
	} else if len(args) < 1 {
		die("Error: not enough arguments for TypeInfo")
	}

	typeName := args[0]
	if typeName == "" {
		die("Error: type is empty")
	}

	if objs, ok := types[typeName]; ok {
		for _, obj := range objs {
			fmt.Println(obj)
		}
	}
}

// Defines access rights.
// Domain name and type must NOT be non-empty strings
// If the domain name OR type name does not exist, they will be created.
// If the operation already exists for that domain and type, it will not
// be added and silently fail.
func addAccess(args []string) {
	if len(args) > 3 {
		die("Error: too many arguments for AddAccess")
	} else if len(args) < 3 {
		die("Error: missing operands for AddAccess")
	}

	op, dName, tName := args[0], args[1], args[2]
	if op == "" {
		die("Error: missing operation")
	} else if dName == "" {
		die("Error missing domain")
	} else if tName == "" {
		die("Error missing type")
	}

	domain, ok := domains[dName]
	if !ok {
		domains[dName] = Domain{}
		domain = Domain{}
	}
	if _, ok := types[tName]; !ok {
		types[tName] = make([]string, 0)
	}

	newPerm := true
	for _, perm := range domain.Permissions {
		if perm.Operation == op && perm.Type == tName {
			// This access permission already exists, don't add it again
			newPerm = false
			break
		}
	}

	if newPerm {
		domain.Permissions = append(domain.Permissions,
			Permission{
				Operation: op,
				Type: tName,
			},
		)
		domains[dName] = domain
	}
	fmt.Println("Success")
}

// Test whether a user can perform an operation on an object
// args[0]: Operation
// arsg[1]: Username
// args[2]: Object
func canAccess(args []string) {
	if len(args) > 3 {
		die("Error: too many arguments for CanAccess")
	} else if len(args) < 3 {
		die("Error: missing operands for CanAccess")
	}

	var user User
	op, userName, obj := args[0], args[1], args[2]

	for _, u := range users {
		if u.Username == userName {
			user = u
		}
	}

	for _, domain := range user.Domains {
		dom, ok := domains[domain]
		if !ok {
			continue
		}
		for _, perm := range dom.Permissions {
			if perm.Operation != op {
				continue
			}
			for _, o := range types[perm.Type] {
				if obj == o {
					fmt.Println("Success")
					return
				}
			}
		}
	}
	die("Error: access denied")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		die("Error: no command given")
	}

	fetchTables()
	switch command, cargs := args[0], args[1:]; command {
	case "AddUser":
		addUser(cargs)
	case "Authenticate":
		authenticate(cargs)
	case "SetDomain":
		setDomain(cargs)
	case "DomainInfo":
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
		die("Error: invalid command %s", command)

	}
	sync()
	os.Exit(0)
}
