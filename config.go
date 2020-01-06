package extract

import (
	"flag"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path/filepath"
)

type HashOptions struct {
	From plumbing.Hash
	To   plumbing.Hash
}

type ConfigOptions struct {
	TargetDir string
	Debug     bool
	Hash      *HashOptions
}

var Config *ConfigOptions

func init() {
	var from, to, targetDir string
	var debug bool
	flag.StringVar(&from, "from", "", "git log begin hash")
	flag.StringVar(&to, "to", "", "git log end hash")
	flag.StringVar(&targetDir, "out", getExecutePath(), "git change output dir")
	flag.BoolVar(&debug, "debug", true, "open debug log")
	flag.Parse()

	var fromHash plumbing.Hash
	if from != "" {
		fromHash = plumbing.NewHash(from)
	}

	var toHash plumbing.Hash
	if to != "" {
		toHash = plumbing.NewHash(to)
	}

	Config = &ConfigOptions{
		TargetDir: targetDir,
		Debug:     debug,
		Hash: &HashOptions{
			From: fromHash,
			To:   toHash,
		},
	}
}

func getExecutePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	CheckIfError(err)
	return dir
}
