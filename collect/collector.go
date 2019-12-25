package collect

import (
	"container/list"
	"fmt"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

type Collector struct {
	DeletedList   *list.List
	AmbiguityList *list.List
}

func New() *Collector {
	collector := &Collector{
		DeletedList:   list.New(),
		AmbiguityList: list.New(),
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

func (c *Collector) Display() {
	fmt.Printf("Deleted Files：%d\n", c.DeletedList.Len())
	loopList(c.DeletedList, func(file string) {
		fmt.Printf("\t%s\n", file)
	})
	fmt.Printf("Ambiguity Files：%d\n", c.AmbiguityList.Len())
	loopList(c.AmbiguityList, func(file string) {
		fmt.Printf("\t%s\n", file)
	})
}

func loopList(list *list.List, out func(string)) {
	if list.Len() == 0 {
		return
	}
	for {
		e := list.Front()
		if e == nil {
			break
		}
		out(e.Value.(string))
		list.Remove(e)
	}
}

func appendDeletedFile(c *Collector, fileName string) {
	c.DeletedList.PushBack(fileName)
}

func appendAmbiguityFile(c *Collector, action merkletrie.Action, fileName string) {
	if action == merkletrie.Modify {
		c.AmbiguityList.PushBack(fmt.Sprintf("[M] %s", fileName))
	} else if action == merkletrie.Insert {
		c.AmbiguityList.PushBack(fmt.Sprintf("[I] %s", fileName))
	} else {
		c.AmbiguityList.PushBack(fmt.Sprintf("[D] %s", fileName))
	}
}
