package extract

import (
	"github.com/HarryCU/git-extract/collect"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
	"io"
	"os"
	"path"
)

func (c Changes) Copy(sourceDir, targetDir string, collector *collect.Collector) {
	if !exists(targetDir) {
		return
	}

	actions := c.Actions()
	for _, action := range actions {
		for _, file := range action.Files {
			srcFile := path.Join(sourceDir, file)

			if action.Key == merkletrie.Delete || !exists(srcFile) {
				collector.Append(action.Key, file)
			}

			destFile := path.Join(targetDir, action.Key.String(), file)
			if !exists(srcFile) || exists(destFile) {
				continue
			}

			err := os.MkdirAll(path.Dir(destFile), os.ModeDir)
			CheckIfError(err)

			err = copyFile(srcFile, destFile)
			CheckIfError(err)
		}
	}
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func copyFile(srcFile, destFile string) error {
	fileRead, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer fileRead.Close()
	fileWrite, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer fileWrite.Close()

	buf := make([]byte, 4096)
	for {
		n, err := fileRead.Read(buf)
		if err != nil && err == io.EOF {
			return nil
		}
		_, _ = fileWrite.Write(buf[:n])
	}
}
