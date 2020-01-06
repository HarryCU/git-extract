package collect

import (
	"github.com/HarryCU/git-extract"
	"github.com/HarryCU/git-extract/log"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

type Action struct {
	Key   merkletrie.Action
	Files []string
}

type ActionCollector struct {
	actionMap     map[string]merkletrie.Action
	actionCounter map[merkletrie.Action]int
}

func (ac *ActionCollector) Include(value interface{}) bool {
	change := value.(*extract.Change)
	log.Debug("Processing：(%s) - (%s)", change.A(), change.B())
	buildActionByChanges(change, ac.actionMap, ac.actionCounter)
	return true
}

func (ac *ActionCollector) Build() []*Action {
	result := make([]*Action, 3)
	actionFileMap := make(map[merkletrie.Action]*Action)
	idx := 0
	for action, count := range ac.actionCounter {
		actionFile := &Action{
			Key:   action,
			Files: make([]string, count),
		}
		ac.actionCounter[action] = 0
		result[idx] = actionFile
		actionFileMap[action] = actionFile
		idx++
	}

	for file, action := range ac.actionMap {
		actionFileMap[action].Files[ac.actionCounter[action]] = file
		ac.actionCounter[action] = ac.actionCounter[action] + 1
	}
	return result[0:idx]
}

func buildActionByChanges(change *extract.Change, actionMap map[string]merkletrie.Action, actionCounter map[merkletrie.Action]int) {
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
