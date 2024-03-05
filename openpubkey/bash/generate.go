package main

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
)

type UsersDetails struct {
	Email		string		`yaml:"email"`
	Principals	[]string	
}

type Users struct {
	Users []UsersDetails   `yaml:"user"` 
}

func main(){
	f, err := os.Create("policy.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("-------Writing config-------")
	input := os.Args[1]
	fmt.Println(input)
	user := UsersDetails{
		Email: input,
		Principals: []string{"root"},
	}
	s1 := Users{
		Users : []UsersDetails{user},
	}
	yamlData2, err := yaml.Marshal(&s1)
	if err!=nil{
		fmt.Println("Error while Marshalling. %v", err)
	}
	l, err := f.WriteString(string(yamlData2))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil{
		fmt.Println(err)
		return
	}
}

