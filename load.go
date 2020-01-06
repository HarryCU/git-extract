package extract

import (
	"github.com/HarryCU/git-extract/filter"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

func Load(path string, chain *filter.Chain) map[*object.Commit]Changes {

	r, err := git.PlainOpen(path)
	CheckIfError(err)

	// ... retrieving the HEAD reference
	ref, err := r.Head()
	CheckIfError(err)

	from := Config.Hash.From
	if from.IsZero() {
		from = ref.Hash()
	}
	to := Config.Hash.To

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: from, Order: git.LogOrderCommitterTime})
	CheckIfError(err)

	// ... just iterates over the commits
	changeMap := make(map[*object.Commit]Changes)
	err = cIter.ForEach(func(c *object.Commit) error {
		if !to.IsZero() && c.Hash == to {
			return storer.ErrStop
		}
		changeMap[c] = New(c, chain)
		return nil
	})
	CheckIfError(err)
	return changeMap
}
