package extract

import (
	"fmt"
	"github.com/HarryCU/git-extract/filter"
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

type ChangeFile struct {
	Name   string
	Action merkletrie.Action
}

func New(c *object.Commit, chain *filter.Chain) Changes {
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
		change := &Change{
			a: ca.ID(),
			b: c.ID(),
			c: diff,
		}
		if chain != nil && !chain.Include(change) {
			return nil
		}
		changes[idx] = change
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
