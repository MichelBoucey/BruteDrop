package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {

	var config Config

	var lines []string

	var invalidUser = regexp.MustCompile(`(\D*?\s\d{2}\s\d{2}:\d{2}:\d{2}).*?\sfor\s(invalid\suser\s|)(.+)\sfrom\s(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s`)

	data, err := ioutil.ReadFile("/etc/brutedrop.conf")

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(err)
	}

	// AuthorizedUsers and AuthorizedAddresses can't be both empty
	if len(config.AuthorizedUsers) == 0 && len(config.AuthorizedAddresses) == 0 {
		fmt.Println("You have to add authorized users or IP addresses in /etc/brutedrop.conf")
		os.Exit(0)
	}

	// Get log lines of failed SSH login attempts from journalctl
	out, err := exec.Command("sh", "-c", "journalctl --since \"43830 minutes ago\" -u sshd --no-pager | grep Failed").Output()

	if string(out) == "" {
		fmt.Println("No log lines to process")
		os.Exit(0)
	}

	lines = strings.Split(string(out), "\n")

	// Iterating over log lines
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {

			matches := invalidUser.FindStringSubmatch(lines[i])
			timestamp := "[" + matches[1] + "]"
			origin := matches[3] + "@" + matches[4]

			if isElement(matches[3], config.AuthorizedUsers) {
				logging(config.LogPath, timestamp + " Authorized user " + origin + " failed to login")
			} else if !isElement(matches[3], config.AuthorizedAddresses) {
				logging(config.LogPath, timestamp + " Unauthorized user " + origin + " failed to login")
			}
		}
	}
}

type Config struct {
	Iptables            string   `yaml:"Iptables"`
	IptablesChain       string   `yaml:"IptablesChain"`
	Journalctl          string   `yaml:"Journalctl"`
	LogPath             string   `yaml:"LogPath"`
	AuthorizedUsers     []string `yaml:"AuthorizedUsers"`
	AuthorizedAddresses []string `yaml:"AuthorizedAddresses"`
}

func isElement (e string, l []string) bool {
	for i := 0; i < len(l); i++ {
		if l[i] == e {
			return true
		}
	}
	return false
}

func logging (p string, s string) {
	exec.Command("sh", "-c", "echo " + s + " >> " + p).Run()
}

