package collect

import (
	"fmt"
	"github.com/HarryCU/git-extract/set"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

type Collector struct {
	DeletedList     *set.Set
	AmbiguityList   *set.Set
	actionCollector *ActionCollector
}

type AmbiguityFile struct {
	Action   merkletrie.Action
	FileName string
}

func (c *Collector) Include(value interface{}) bool {
	return c.actionCollector.Include(value)
}

func New() *Collector {
	collector := &Collector{
		DeletedList:   set.New(),
		AmbiguityList: set.New(),
		actionCollector: &ActionCollector{
			actionMap:     make(map[string]merkletrie.Action),
			actionCounter: make(map[merkletrie.Action]int),
		},
	}
	return collector
}

func (c *Collector) Append(action merkletrie.Action, fileName string) {
	switch action {
	case merkletrie.Delete:
		appendDeletedFile(c, fileName)
	default:
		appendAmbiguityFile(c, action, fileName)
	}
}

func (c *Collector) EliminateAmbiguity() {
	cleanSet := set.New()
	loop(c.AmbiguityList, func(value interface{}) {
		file := value.(*AmbiguityFile)
		if c.DeletedList.Contains(file.FileName) {
			_ = cleanSet.AddIfAbsent(file)
		}
	})
	loop(cleanSet, func(value interface{}) {
		file := value.(*AmbiguityFile)
		_ = c.AmbiguityList.Delete(value)
		_ = c.DeletedList.Delete(file.FileName)
	})
}

func (c *Collector) Display() {
	fmt.Printf("Deleted Files：%d\n", c.DeletedList.Size())
	loop(c.DeletedList, func(value interface{}) {
		file := value.(string)
		fmt.Printf("\t%s\n", file)
	})
	fmt.Printf("Ambiguity Files：%d\n", c.AmbiguityList.Size())
	loop(c.AmbiguityList, func(value interface{}) {
		file := value.(*AmbiguityFile)
		fmt.Print("\t[")
		switch file.Action {
		case merkletrie.Insert:
			fmt.Print("I")
		case merkletrie.Modify:
			fmt.Print("M")
		case merkletrie.Delete:
			fmt.Print("D")
		}
		fmt.Printf("] %s\n", file.FileName)
	})
}

func loop(set *set.Set, out func(interface{})) {
	if set.Size() == 0 {
		return
	}
	_, _ = set.ForEach(func(value interface{}) bool {
		out(value)
		return true
	})
}

func appendDeletedFile(c *Collector, fileName string) {
	c.DeletedList.AddIfAbsent(fileName)
}

func appendAmbiguityFile(c *Collector, action merkletrie.Action, fileName string) {
	file := &AmbiguityFile{
		Action:   action,
		FileName: fileName,
	}
	_ = c.AmbiguityList.AddIfAbsent(file)
}
