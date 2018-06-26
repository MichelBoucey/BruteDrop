package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
        "log"
	// import "regexp"
)

type Config struct {
	Iptables            string   `yaml:"Iptables"`
	IptablesChain       string   `yaml:"IptablesChain"`
	Journalctl          string   `yaml:"Journalctl"`
	AuthorizedUsers     []string `yaml:"AuthorizedUsers"`
	AuthorizedAddresses []string `yaml:"AuthorizedAddresses"`
}

func main() {

	var config Config

	data, err := ioutil.ReadFile("conf/brutedrop.conf")

	if err != nil {
                log.Fatal(err)
	}

	// fmt.Println(string(data))

	err = yaml.Unmarshal(data,&config)

	if err != nil {
                log.Fatal(err)
	}

}
