package main

import (
	"fmt"
	extract "github.com/HarryCU/git-extract"
	"os"
)

func main() {
	path := os.Args[len(os.Args)-1]
	changeMap := extract.Load(path)

	for commit, changes := range changeMap {
		fmt.Printf("Commit ID：%s\n", commit.ID())
		fmt.Print("Change Files：\n")

		actions := changes.Actions()
		for _, action := range actions {
			fmt.Printf(" \t%s：\n", action.Key)
			for _, file := range action.Files {
				fmt.Printf(" \t  %s\n", file)
			}
		}
		changes.Copy(path, extract.Config.TargetDir)
		fmt.Print("\n")
	}
}
