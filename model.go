package extract

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

type Changes []*Change

type Change struct {
	a plumbing.Hash
	b plumbing.Hash
	c object.Changes
}

type Action struct {
	Key   merkletrie.Action
	Files []string
}

type ChangeFile struct {
	Name   string
	Action merkletrie.Action
}

func New(c *object.Commit) Changes {
	b, err := c.Tree()
	CheckIfError(err)
	changes := make(Changes, c.NumParents())
	idx := 0
	err = c.Parents().ForEach(func(ca *object.Commit) error {
		a, err := ca.Tree()
		CheckIfError(err)

		diff, err := object.DiffTree(a, b)
		CheckIfError(err)

		if len(diff) == 0 {
			return nil
		}

		changes[idx] = &Change{
			a: ca.ID(),
			b: c.ID(),
			c: diff,
		}
		idx++
		return nil
	})
	CheckIfError(err)
	return changes[0:idx]
}

func (c *Change) A() string {
	return c.a.String()
}

func (c *Change) B() string {
	return c.b.String()
}

func (c *Change) Files() []*ChangeFile {
	files := make([]*ChangeFile, len(c.c))
	for i, changes := range c.c {
		action, err := changes.Action()
		CheckIfError(err)
		switch action {
		case merkletrie.Insert:
			files[i] = &ChangeFile{
				Name:   changes.To.Name,
				Action: action,
			}
		case merkletrie.Delete, merkletrie.Modify:
			files[i] = &ChangeFile{
				Name:   changes.From.Name,
				Action: action,
			}
		default:
			fmt.Printf("unknown action: %d\n", action)
		}
	}
	return files
}

func (c Changes) Actions() []*Action {

	actionMap := make(map[string]merkletrie.Action)
	actionCounter := make(map[merkletrie.Action]int)
	for _, change := range c {
		files := change.Files()
		for _, file := range files {
			action, ok := actionMap[file.Name]
			if !ok {
				actionMap[file.Name] = file.Action
				// 计数
				actionCounter[file.Action] = actionCounter[file.Action] + 1
			} else {
				// 在查找区间，发现其文件有Insert状态，就代表新增文件
				if file.Action == merkletrie.Insert {
					actionMap[file.Name] = merkletrie.Insert
					// 重新计数
					actionCounter[action] = actionCounter[action] - 1
					actionCounter[file.Action] = actionCounter[file.Action] + 1
				}
			}
		}
	}

	result := make([]*Action, 3)
	actionFileMap := make(map[merkletrie.Action]*Action)
	idx := 0
	for action, count := range actionCounter {
		actionFile := &Action{
			Key:   action,
			Files: make([]string, count),
		}
		actionCounter[action] = 0
		result[idx] = actionFile
		actionFileMap[action] = actionFile
		idx++
	}

	for file, action := range actionMap {
		actionFileMap[action].Files[actionCounter[action]] = file
		actionCounter[action] = actionCounter[action] + 1
	}

	return result[0:idx]
}
