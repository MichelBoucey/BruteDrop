package main

import (
	"flag"
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

	version := "1.0.0"

	var config Config

	var logLines []string

	var invalidUser = regexp.MustCompile(`^(.*?\d{2}:\d{2}:\d{2}).*?Invalid\suser\s(\w+)\sfrom\s(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\sport\s\d{1,5}$`)

	versionFlag := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *versionFlag == true {
		fmt.Println("brutedrop v"+version+"\nCopyright © 024 Michel Boucey\nReleased under 3-Clause BSD License")
		os.Exit(0)
	}

	// Get and check bruteDrop configuration
	data, err := ioutil.ReadFile("/etc/brutedrop.conf")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(config.Iptables); os.IsNotExist(err) {
		log.Fatal("Can't find iptables at path "+config.Iptables)
	}
	if _, err := os.Stat(config.Journalctl); os.IsNotExist(err) {
		log.Fatal("Can't find journalctl at path "+config.Journalctl)
	}
	// AuthorizedUsers and AuthorizedAddresses can't be both empty
	if len(config.AuthorizedUsers) == 0 && len(config.AuthorizedAddresses) == 0 {
		log.Fatal("To run brutedrop you have to add authorized users and/or IP addresses in /etc/brutedrop.conf")
	}

	// Get log lines of failed SSH login attempts from journalctl
	out, err := exec.Command("sh", "-c", config.Journalctl+" --since \""+strconv.Itoa(config.LogEntriesSince)+" minutes ago\" -u sshd --no-pager | grep Invalid").Output()

	if len(out) == 0 {
		os.Exit(0)
	}
	logLines = strings.Split(string(out), "\n")

	// Iterating over log lines
	for i := 0; i < len(logLines); i++ {

		if logLines[i] != "" {

			matches := invalidUser.FindStringSubmatch(logLines[i])

			if len(matches) == 4 {

				if isElement(matches[2], config.AuthorizedUsers) {

					logging(config.LoggingTo, "brutedrop: Authorized user "+matches[2]+" failed to login from "+matches[3]+" at "+matches[1])

				} else if !isElement(matches[3], config.AuthorizedAddresses) {

					// Is this IP address already banned ?
					_, err := exec.Command("sh", "-c", config.Iptables+" -w -C INPUT -s "+matches[3]+" -j DROP").Output()
					if err != nil {
						// Ban IP address
						appendRule := config.Iptables+" -w -A INPUT -s "+matches[3]+" -j DROP"
						err := exec.Command("sh", "-c", appendRule).Run()
						if err != nil {
							log.Fatal("Can't execute \""+appendRule+"\"")
						}
						logging(config.LoggingTo, "brutedrop: Dropping "+matches[3]+" from invalid user "+matches[2]+" connection at "+matches[1])
					}
				} else {

					logging(config.LoggingTo, "brutedrop: Invalid user "+matches[2]+" from authorized IP address "+matches[3])

				}
			}
		}
	}
}

type Config struct {
	Iptables            string   `yaml:"Iptables"`
	Journalctl          string   `yaml:"Journalctl"`
	LoggingTo           string   `yaml:"LoggingTo"`
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
