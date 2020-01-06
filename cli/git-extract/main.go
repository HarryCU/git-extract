package main

import (
	"fmt"
	extract "github.com/HarryCU/git-extract"
	"github.com/HarryCU/git-extract/collect"
	"github.com/HarryCU/git-extract/filter"
	"github.com/HarryCU/git-extract/log"
	"os"
)

func init() {
	log.Init()
}

func main() {
	sourceDir := os.Args[len(os.Args)-1]

	collector := collect.New()
	chain := filter.New(collector)
	changeMap := extract.Load(sourceDir, chain)
	fmt.Print("Commit ID(s)：\n")
	for commit := range changeMap {
		fmt.Printf("\t%s\n", commit.ID())
	}
	collector.CopyTo(sourceDir, extract.Config.TargetDir)
	collector.EliminateAmbiguity()
	collector.Display()
}
