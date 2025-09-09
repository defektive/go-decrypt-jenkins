package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/thesubtlety/go-decrypt-jenkins/pkg/config"
	"github.com/thesubtlety/go-decrypt-jenkins/pkg/jenkinscrypto"
	"github.com/thesubtlety/go-decrypt-jenkins/pkg/worker"
)

func usage() {
	fmt.Printf("Usage: %s -m master.key -s hudson.util.Secret -c credentials.xml\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.StringVar(&config.Masterkeyfile, "m", "", "master.key")
	flag.StringVar(&config.Secretkeyfile, "s", "", "hudson.util.Secret")
	flag.StringVar(&config.Secretbytesfile, "sb", "", "com.cloudbees.plugins.credentials.SecretBytes.KEY")
	flag.StringVar(&config.Credfile, "c", "", "credentials.xml")
	flag.StringVar(&config.Searchdirectory, "d", "", "directory to search for xml files, and appropiate key files")
	flag.StringVar(&config.PluginExperimental, "p", "", "plugin to search for, e.g. com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl")
	flag.BoolVar(&config.Optionbrute, "b", false, "decrypt everything that looks {encrypty}")
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) == 1 {
		usage()
	}

	//append xml node name to list
	if config.PluginExperimental != "" {
		worker.Plugins = append(worker.Plugins, config.PluginExperimental)
	}

	//search a directory for all the key files and xml files to parse and decrypt
	if config.Searchdirectory != "" {
		worker.Search()
		return
	}

	//exit if missing master, secret, credfile
	if config.Masterkeyfile == "" || config.Credfile == "" && (config.Secretbytesfile == "" || config.Secretkeyfile == "") {
		usage()
	}

	//regex decrypt single file with secretbytes
	if config.Secretbytesfile != "" {
		sbk := jenkinscrypto.Initsecretbytesdecrypt(config.Secretbytesfile, config.Masterkeyfile)
		worker.Brute(sbk)
		return
	}

	//regex decrypt single file
	if config.Optionbrute {
		k := jenkinscrypto.Initdecrypt(config.Secretkeyfile, config.Masterkeyfile)
		worker.Brute(k)
		return
	}

	// default, try to decrypt the specified credfile given the cli options
	secrets := worker.DefaultParse()

	b, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile("secrets.json", b, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
