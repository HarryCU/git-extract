package extract

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func NewDiff(sourceDir string) (ca, cb *object.Commit, diff object.Changes) {

	r, err := git.PlainOpen(sourceDir)
	CheckIfError(err)

	cb, err = r.CommitObject(Config.Hash.From)
	CheckIfError(err)
	ca, err = r.CommitObject(Config.Hash.To)
	CheckIfError(err)

	a, err := ca.Tree()
	CheckIfError(err)
	b, err := cb.Tree()
	CheckIfError(err)

	diff, err = object.DiffTree(a, b)
	CheckIfError(err)
	return
}
