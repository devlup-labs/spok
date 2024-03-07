package policy

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"

	"gopkg.in/yaml.v2"
)

type UsersDetails struct {
	Email      string   `yaml:"email"`
	Principals []string `yaml:"principals"`
}

type Policy struct {
	User []UsersDetails `yaml:"user"`
}

func getFileSize(file string) (int64, error) {
	fileinfo, err := os.Stat(file)
	if err != nil {
		return 0, err
	}

	return fileinfo.Size(), nil
}

func (p *Policy) unMarshal(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// var details []UsersDetails
	maxCapacity, _ := getFileSize(filename)
	reader := bufio.NewReader(file)
	data := make([]byte, maxCapacity)
	_, err = reader.Read(data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// p.Users = details
	return yaml.Unmarshal(data, &p)
}

func AddPolicy(email string, principal string) {
	newInstance := new(Policy)
	err := newInstance.unMarshal("../policy.yml")
	if err != nil {
		log.Fatalln(err)
	}
	i := 0
	idx := 0
	toAdd := new(UsersDetails)
	toAdd.Email = email
	toAdd.Principals = []string{principal}
	for _, u := range newInstance.User {
		if u.Email == email {
			if !slices.Contains(u.Principals, principal) {
				u.Principals = slices.Concat(u.Principals, toAdd.Principals)
				fmt.Println("Email exists adding Permission")
				newInstance.User[idx].Principals = u.Principals
				i = 1
				break
			} else {
				i = 1
				fmt.Println("User Exists")
				break
			}
		}
		idx += 1
	}
	if i != 1 {
		newInstance.User = append(newInstance.User, *toAdd)
	}
	fmt.Println(newInstance)
	f, err := os.Create("../policy.yml")
	if err != nil {
		fmt.Println("File not Found policy.yml")
	}
	yamlData2, err := yaml.Marshal(&newInstance)
	if err != nil {
		fmt.Println("Error while Marshalling. ", err)
	}
	l, err := f.WriteString(string(yamlData2))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
