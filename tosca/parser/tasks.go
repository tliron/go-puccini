package parser

import (
	"sort"
	"strings"

	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca/parsing"
)

type Executor func(task *Task)

//
// Task
//

type Task struct {
	Name         string
	Executor     Executor
	Parents      Tasks
	Dependencies Tasks
}

type Tasks map[*Task]struct{}

type TasksForEntities map[parsing.EntityPtr]*Task

func NewTask(name string) *Task {
	return &Task{Name: name, Parents: make(Tasks), Dependencies: make(Tasks)}
}

func (self *Task) IsIndependent() bool {
	return len(self.Dependencies) == 0
}

func (self *Task) Execute() {
	// If we got here, we should be independent (no dependencies)
	self.Executor(self)
}

func (self *Task) Done() {
	// If we got here, we should be independent (no dependencies)
	for parent := range self.Parents {
		parent.Dependencies.Remove(self)
	}

	// Make sure we won't be reused
	self.Executor = nil
	self.Parents = nil
	self.Dependencies = nil
}

func (self *Task) Print(indent int) {
	terminal.PrintIndent(indent)
	terminal.Printf("%s\n", terminal.StdoutStylist.Path(self.Name))
	self.PrintDependencies(indent, terminal.TreePrefix{})
}

func (self *Task) PrintDependencies(indent int, treePrefix terminal.TreePrefix) {
	// Sort
	var taskList TaskList
	for dependency := range self.Dependencies {
		taskList = append(taskList, dependency)
	}
	sort.Sort(taskList)

	last := len(taskList) - 1
	for i, dependency := range taskList {
		isLast := i == last
		dependency.PrintDependency(indent, treePrefix, isLast)
		dependency.PrintDependencies(indent, append(treePrefix, isLast))
	}
}

func (self *Task) PrintDependency(indent int, treePrefix terminal.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	terminal.Printf("%s\n", self.Name)
}

func (self *Task) AddDependency(task *Task) {
	self.Dependencies.Add(task)
	task.Parents.Add(self)
}

//
// Tasks
//

func (self Tasks) Add(task *Task) {
	self[task] = struct{}{}
}

func (self Tasks) Remove(task *Task) {
	delete(self, task)
}

func (self Tasks) FindIndependent() (*Task, bool) {
	for task := range self {
		if task.IsIndependent() {
			return task, true
		}
	}
	return nil, false
}

func (self Tasks) Validate() bool {
	// TODO make sure there are no endless loops
	return true
}

func (self Tasks) Drain() {
	if !self.Validate() {
		return
	}

	logTasks.Debugf("starting %d tasks", len(self))

	for true {
		task, ok := self.FindIndependent()
		if !ok {
			break
		}

		self.Remove(task)
		task.Execute()

		// After one independent task is done, other tasks should become independent
	}

	if len(self) > 0 {
		logTasks.Warningf("%d tasks not completed", len(self))
	}
}

// Print

func (self Tasks) Print(indent int) {
	// Sort
	var taskList TaskList
	for task := range self {
		taskList = append(taskList, task)
	}
	sort.Sort(taskList)

	for _, task := range taskList {
		task.Print(indent)
	}
}

type TaskList []*Task

// ([sort.Interface])
func (self TaskList) Len() int {
	return len(self)
}

// ([sort.Interface])
func (self TaskList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// ([sort.Interface])
func (self TaskList) Less(i, j int) bool {
	return strings.Compare(self[i].Name, self[j].Name) < 0
}
