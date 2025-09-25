// Package flags
package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type FlagConfig struct {
	User   string
	Repo   string
	Branch string
	Args   []string
}

func Load() (*FlagConfig, error) {
	fcfg := &FlagConfig{}

	flag.StringVar(&fcfg.User, "name", "", "Github user name")
	flag.StringVar(&fcfg.User, "n", "", "Github user name")

	flag.StringVar(&fcfg.Repo, "repo", "", "Github repo name")
	flag.StringVar(&fcfg.Repo, "r", "", "Github repo name")

	flag.StringVar(&fcfg.Branch, "b", "", "Github branch name")
	flag.StringVar(&fcfg.Branch, "branch", "", "Github branch name")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: myapp [options]")
		fmt.Fprintf(os.Stderr, "\nRequired flags:")
		fmt.Fprintf(os.Stderr, "  -n, -name     GitHub user name")
		fmt.Fprintf(os.Stderr, "  -r, -repo     GitHub repo name")
		fmt.Fprintf(os.Stderr, "  -b, -branch   GitHub branch name")
		fmt.Fprintf(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()
	fcfg.Args = flag.Args()

	if err := fcfg.CheckFlag(); err != nil {
		return nil, err
	}

	if err := fcfg.CheckArgs(); err != nil {
		return nil, err
	}

	return fcfg, nil
}

func (cfg *FlagConfig) CheckArgs() error {
	if len(cfg.Args) == 0 {
		return errors.New("missing file name")
	}
	return nil
}

func (cfg *FlagConfig) CheckFlag() error {
	var missing []string

	if cfg.User == "" {
		missing = append(missing, "user")
	}
	if cfg.Branch == "" {
		missing = append(missing, "branch")
	}
	if cfg.Repo == "" {
		missing = append(missing, "repo")
	}

	if len(missing) > 0 {
		return errors.New("missing required flags: " + strings.Join(missing, ", "))
	}
	return nil
}

func (cfg *FlagConfig) String() string {
	format := `Username: %s\nRepo: %s\nBranch: %s\nArgs:%v\n`
	return fmt.Sprintf(format,
		cfg.User,
		cfg.Repo,
		cfg.Branch,
		cfg.Args,
	)
}
