package main

import (
	"fmt"
	extract "github.com/HarryCU/git-extract"
	"github.com/HarryCU/git-extract/collect"
	"os"
)

func main() {
	path := os.Args[len(os.Args)-1]
	changeMap := extract.Load(path)

	collector := collect.New()

	fmt.Print("Commit ID(s)：\n")
	for commit, changes := range changeMap {
		fmt.Printf("\t%s\n", commit.ID())
		changes.Copy(path, extract.Config.TargetDir, collector)
	}

	collector.Display()
}
