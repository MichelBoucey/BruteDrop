package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
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

	// var invalidUser = regexp.MustCompile(`/^(\D{3}\s\d{2}\s\d{2}:\d{2}:\d{2}).*?for\s(invalid user\s|)(.+)\sfrom\s(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/`)
	var invalidUser = regexp.MustCompile(`\sfor\s(.+)\sfrom\s(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s`)

	data, err := ioutil.ReadFile("conf/brutedrop.conf")

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(err)
	}

	// Get SSH log lines from journaltctl
	// out, err := exec.Command("bash","-c","journalctl --since \"5 minutes ago\" -u sshd --no-pager | grep Failed").Output()
	out, err := exec.Command("bash", "-c", "journalctl --since \"1000 minutes ago\" -u sshd --no-pager | grep Accepted").Output()

	if string(out) == "" {
		fmt.Println("No log lines to process")
	}

	lines = strings.Split(string(out), "\n")

	// Iterating over lines
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {
			matches := invalidUser.FindStringSubmatch(lines[i])
			fmt.Println(matches[2])
		}
	}
}
