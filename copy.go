package extract

import (
	"io"
	"os"
	"path"
)

func (c Changes) Copy(sourceDir, targetDir string) bool {
	if !exists(targetDir) {
		return false
	}

	actions := c.Actions()
	for _, action := range actions {
		for _, file := range action.Files {
			srcFile := path.Join(sourceDir, file)
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
	return true
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
