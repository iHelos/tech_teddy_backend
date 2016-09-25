package deploy_config

import (

	"encoding/json"
	"log"
	"fmt"
	"io/ioutil"
	"os"
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
	file, e := ioutil.ReadFile(filepath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	//m := new(Dispatch)
	//var m interface{}

	configuration := DeployConfiguration{}
	err := json.Unmarshal(file, &configuration)
	if err!=nil{
		log.Print(err)
	}
	return &configuration
}