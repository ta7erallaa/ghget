// Package main
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ta7eralla/ghget/config"
	"github.com/ta7eralla/ghget/downloader"
	"github.com/ta7eralla/ghget/flags"
)

// TODO:
// Support more files.
// Fix typos in function names
// Support one link like wget
// add flag to specify location of file

func main() {
	log.SetFlags(0)

	flagConfig, err := flags.Load()
	if err != nil {
		log.Fatal(err)
	}

	if flagConfig.IsOneFlagSet() {
		if err := flagConfig.ValidateFlags(); err != nil {
			log.Fatal(err)
		}
	}

	fileConfig, err := getConfig(flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	filenames := flagConfig.Filenames

	downloader := downloader.New()
	if err := downloader.DownloadFromConfig(fileConfig, filenames); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func getConfig(fcfg *flags.FlagConfig) (*config.Config, error) {
	if fcfg.IsAllFlagSet() {
		cfg := config.New(fcfg.Name, fcfg.Repo, fcfg.Branch)

		if !cfg.IsNewFlagsEqualConfigValues() {
			if err := cfg.Write(); err != nil {
				fmt.Fprint(os.Stderr, "failed to save new flag to config.json file")
			}
		}
		return cfg, nil
	}

	cfg := &config.Config{}
	if err := cfg.Read(); err != nil {
		log.SetPrefix("config path: ")
		return nil, err
	}
	return cfg, nil
}

