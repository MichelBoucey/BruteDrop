package main

type Config struct {
	DryRun              bool     `yaml:"DryRunMode"`
	Iptables            string   `yaml:"IptablesBinPath"`
	Journalctl          string   `yaml:"JournalctlBinPath"`
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
