package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
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

	var lines []string

	data, err := ioutil.ReadFile("conf/brutedrop.conf")

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(err)
	}

	// out, err := exec.Command("bash","-c","journalctl --since \"5 minutes ago\" -u sshd --no-pager | grep Failed").Output()
	out, err := exec.Command("bash", "-c", "journalctl --since \"5 minutes ago\" -u sshd --no-pager").Output()

	if string(out) == "" {
		fmt.Println("Nil")
	}

	lines = strings.Split(string(out), "\n")

	// Iterating over lines
	for i := 0; i < len(lines); i++ {
		// /^(\D{3}\s\d{2}\s\d{2}:\d{2}:\d{2}).*?for\s(invalid user\s|)(.+)\sfrom\s(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/
		fmt.Print("> ")
		fmt.Println(lines[i])
		fmt.Println()

	}
}
