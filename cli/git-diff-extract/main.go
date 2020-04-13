package main

import (
	extract "github.com/HarryCU/git-extract"
	"github.com/HarryCU/git-extract/collect"
	"github.com/HarryCU/git-extract/log"
	"os"
)

func init() {
	log.Init()
}

func main() {
	sourceDir := os.Args[len(os.Args)-1]

	ca, cb, diff := extract.NewDiff(sourceDir)

	change := extract.NewSingle(ca, cb, diff)
	collector := collect.New()
	collector.Include(change)

	collector.CopyTo(sourceDir, extract.Config.TargetDir)
	collector.EliminateAmbiguity()
	collector.Display()
}
