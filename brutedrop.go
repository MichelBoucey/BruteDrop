package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	// import "regexp"
)

type Config struct {
	Iptables            string   //   `yaml:"iptables"`
	IptablesChain       string   //   `yaml:"iptables_chain"`
	Journalctl          string   //   `yaml:"journalctl"`
	AuthorizedUsers     []string //   `yaml:"authorized_users"`
	AuthorizedAddresses []string //   `yaml:"authorized_addresses"`
}

func main() {

	var config Config

	data, rferr := ioutil.ReadFile("btedrop.conf")

	if rferr != nil {
		fmt.Errorf("Error reading configuration file: %v", rferr)
	}

	umerr := yaml.UnmarshalStrict(data, &config)

	if umerr != nil {
		fmt.Errorf("Configuration file errors: %v", umerr)
	}

	// fmt.Println("bruteDrop")
	fmt.Println(config.Iptables)

}
