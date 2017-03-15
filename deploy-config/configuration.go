package deploy_config

import (

	"encoding/json"
	"log"
	"fmt"
	"io/ioutil"
)

type DeployConfiguration struct {
	Host    string	`json:"host"`
	Port 	string	`json:"port"`
	Database 	DatabaseConfiguration	`json:"database"`
}

type DatabaseConfiguration struct {
	Host string	`json:"host"`
	User string	`json:"user"`
	Password string	`json:"password"`
}

func GetConfiguration (filepath string) *DeployConfiguration{
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read configuration: %s", err.Error())
	}
	fmt.Printf("%s\n", string(file))

	//m := new(Dispatch)
	//var m interface{}

	configuration := DeployConfiguration{}
	err = json.Unmarshal(file, &configuration)
	if err!=nil{
		log.Print(err)
	}
	return &configuration
}