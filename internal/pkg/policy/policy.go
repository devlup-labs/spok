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

func getFileSize(filename string) (int64, error) {
	fileinfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	return fileinfo.Size(), nil
}

func (p *Policy) Unmarshal(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileSize, _ := getFileSize(filename)
	data := make([]byte, fileSize)

	reader := bufio.NewReader(file)

	_, err = reader.Read(data)
	if err != nil {
		fmt.Println(err)

		return err
	}

	return yaml.Unmarshal(data, &p)
}

func AddPolicy(email string, principal string) {
	policy := new(Policy)

	err := policy.Unmarshal("/etc/sos/policy.yml")
	if err != nil {
		log.Fatalln(err)
	}

	i := 0
	idx := 0

	userDetails := new(UsersDetails)
	userDetails.Email = email
	userDetails.Principals = []string{principal}

	for _, u := range policy.User {
		if u.Email == email {
			if !slices.Contains(u.Principals, principal) {
				u.Principals = slices.Concat(
					u.Principals, userDetails.Principals,
				)

				fmt.Println("Email exists adding Permission")

				policy.User[idx].Principals = u.Principals

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
		policy.User = append(policy.User, *userDetails)
	}

	f, err := os.Create("/etc/sos/policy.yml")
	if err != nil {
		fmt.Println("File not Found policy.yml")

		return
	}
	defer f.Close()

	yamlData, err := yaml.Marshal(&policy)
	if err != nil {
		fmt.Println("Error while Marshalling. ", err)

		return
	}

	_, err = f.WriteString(string(yamlData))
	if err != nil {
		fmt.Println(err)

		return
	}
}
