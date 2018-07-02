package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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

	if _, err := os.Stat(config.Iptables); os.IsNotExist(err) {
		log.Fatal("Can't find iptables at path " + config.Iptables)
	}

	if _, err := os.Stat(config.Journalctl); os.IsNotExist(err) {
		log.Fatal("Can't find journalctl at path " + config.Journalctl)
	}

	// AuthorizedUsers and AuthorizedAddresses can't be both empty
	if len(config.AuthorizedUsers) == 0 && len(config.AuthorizedAddresses) == 0 {
		log.Fatal("To run brutedrop you have to add authorized users and/or IP addresses in /etc/brutedrop.conf")
	}

	// Get log lines of failed SSH login attempts from journalctl
	out, err := exec.Command("sh", "-c", config.Journalctl+" --since \""+strconv.Itoa(config.LogEntriesSince)+" minutes ago\" -u sshd --no-pager | grep Failed").Output()

	if len(out) == 0 {
		os.Exit(0)
	}

	lines = strings.Split(string(out), "\n")

	// Iterating over log lines
	for i := 0; i < len(lines); i++ {
		if lines[i] != "" {

			matches := invalidUser.FindStringSubmatch(lines[i])
			timestamp := "[" + matches[1] + "]"

			if isElement(matches[3], config.AuthorizedUsers) {
				logging(config.Logging, timestamp+" Authorized user "+matches[3]+" failed to login from"+matches[4])
			} else if !isElement(matches[4], config.AuthorizedAddresses) {
				_, err := exec.Command("sh", "-c", config.Iptables+" -w -C INPUT -s "+matches[4]+" -j DROP").Output()
				if err != nil {
					appendRule := config.Iptables + " -w -A INPUT -s " + matches[4] + " -j DROP"
					err := exec.Command("sh", "-c", appendRule).Run()
					if err != nil {
						log.Fatal("Can't execute \"" + appendRule + "\"")
					}
					logging(config.Logging, timestamp+" DROP "+matches[3]+" from "+matches[4])
				}
			}
		}
	}
}

type Config struct {
	Iptables            string   `yaml:"Iptables"`
	Journalctl          string   `yaml:"Journalctl"`
	Logging             string   `yaml:"Logging"`
	LogEntriesSince     int      `yaml:"LogEntriesSince"`
	AuthorizedUsers     []string `yaml:"AuthorizedUsers"`
	AuthorizedAddresses []string `yaml:"AuthorizedAddresses"`
}

func isElement(e string, l []string) bool {
	for i := 0; i < len(l); i++ {
		if l[i] == e {
			return true
		}
	}
	return false
}

func logging(p string, s string) {
	if p != "stdout" {
		exec.Command("sh", "-c", "echo "+s+" >> "+p).Run()
	} else {
		fmt.Println(s)
	}
}
