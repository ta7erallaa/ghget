// Package flags
package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type FlagConfig struct {
	Name      string
	Repo      string
	Branch    string
	Filenames []string
}

func Load() (*FlagConfig, error) {
	fcfg := &FlagConfig{}

	flag.StringVar(&fcfg.Name, "name", "", "Github user name")
	flag.StringVar(&fcfg.Name, "n", "", "Github user name")

	flag.StringVar(&fcfg.Repo, "repo", "", "Github repo name")
	flag.StringVar(&fcfg.Repo, "r", "", "Github repo name")

	flag.StringVar(&fcfg.Branch, "b", "", "Github branch name")
	flag.StringVar(&fcfg.Branch, "branch", "", "Github branch name")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: ghget [options]")
		fmt.Fprintf(os.Stderr, "\nRequired flags:\n")
		fmt.Fprintf(os.Stderr, "  -n, -name     GitHub user name\n")
		fmt.Fprintf(os.Stderr, "  -r, -repo     GitHub repo name\n")
		fmt.Fprintf(os.Stderr, "  -b, -branch   GitHub branch name\n")
		fmt.Fprintf(os.Stderr, "\n\nOptions:\n")
	}

	flag.Parse()

	fcfg.Filenames = flag.Args()

	if err := fcfg.ValidateFilename(); err != nil {
		return nil, err
	}

	return fcfg, nil
}

func (cfg *FlagConfig) ValidateFilename() error {
	if len(cfg.Filenames) == 0 {
		return errors.New("missing file name")
	}

	if len(cfg.Filenames) > 1 {
		return errors.New("only accept one file for right now")
	}

	return nil
}

func (cfg *FlagConfig) ValidateFlags() error {
	if cfg.Name == "" {
		return errors.New("(n)ame flag is empty")
	}

	if cfg.Branch == "" {
		return errors.New("(b)ranch flag is empty")
	}

	if cfg.Repo == "" {
		return errors.New("(r)epo flag is empty")
	}

	return nil
}

func (cfg *FlagConfig) IsOneFlagSet() bool {
	return cfg.Name != "" || cfg.Branch != "" || cfg.Repo != ""
}

func (cfg *FlagConfig) IsAllFlagSet() bool {
	return cfg.Name != "" && cfg.Branch != "" && cfg.Repo != ""
}

func (cfg *FlagConfig) String() string {
	format := "Username: %s\nRepo: %s\nBranch: %s\nFilename: %v\n"
	return fmt.Sprintf(format,
		cfg.Name,
		cfg.Repo,
		cfg.Branch,
		cfg.Filenames,
	)
}
