package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

/*
 * Features to add
 *
 * 1. Deadlines
 *    (i)   Give each task an optional deadline with the 'deadline <ID> <date>' command.
 *    (ii)  Show all tasks due today with the 'today' command.
 * 2. Customisable IDs
 *    (i)   Allow the user to specify an identifier that's not a number.
 *    (ii)  Disallow spaces and special characters from the ID.
 * 3. Deletion
 *    (i)   Allow users to delete tasks with the 'delete <ID>' command.
 * 4. Views
 *    (i)   View tasks by date with the 'view by date' command.
 *    (ii)  View tasks by deadline with the 'view by deadline' command.
 *    (iii) Don't remove the functionality that allows users to view tasks by project,
 *          but change the command to 'view by project'
 */

const (
	// Quit is the text command used to quit the task manager.
	Quit   string = "quit"
	prompt string = "> "
)

// TaskList is a set of tasks, grouped by project.
type TaskList struct {
	in  io.Reader
	out io.Writer

	projectTasks Projects
	lastID       TaskId
	timeGetter   TimeGetter
}

// NewTaskList initializes a TaskList on the given I/O descriptors.
func NewTaskList(in io.Reader, out io.Writer) *TaskList {
	return &TaskList{
		in:           in,
		out:          out,
		projectTasks: make(Projects),
		lastID:       0,
		timeGetter:   time.Now,
	}
}

type TimeGetter func() time.Time

func (l *TaskList) WithCalendar(getter TimeGetter) *TaskList {
	l.timeGetter = getter
	return l
}

// Run runs the command loop of the task manager.
// Sequentially executes any given command, until the user types the Quit message.
func (l *TaskList) Run() {
	scanner := bufio.NewScanner(l.in)

	fmt.Fprint(l.out, prompt)
	for scanner.Scan() {
		cmdLine := scanner.Text()
		if cmdLine == Quit {
			return
		}

		l.execute(cmdLine)
		fmt.Fprint(l.out, prompt)
	}
}

type Command string

const (
	Show Command = "show"
)

func (l *TaskList) execute(cmdLine string) {
	args := strings.Split(cmdLine, " ")
	command := Command(args[0])
	switch command {
	case Show:
		l.show()
	case "add":
		l.add(args[1:])
	case "check":
		l.check(args[1])
	case "uncheck":
		l.uncheck(args[1])
	case "deadline":
		taskID, _ := newTaskIdFromString(args[1])
		deadline, _ := NewDeadlineFromString(args[2])
		l.deadline(taskID, deadline)
	case "today":
		l.today()
	case "help":
		l.help()
	default:
		l.error(command)
	}
}

func (l *TaskList) help() {
	fmt.Fprintln(l.out, `Commands:
  show
  add project <project name>
  add task <project name> <task description>
  check <task ID>
  uncheck <task ID>
  `)
}

func (l *TaskList) error(command Command) {
	fmt.Fprintf(l.out, "Unknown command \"%s\".\n", command)
}

func (l *TaskList) show() {
	// show projects sequentially
	l.printSortedProjects(l.projectTasks.SortedProjects(), anyTask)
}

type Category string

const (
	CategoryProject Category = "project"
	CategoryTask    Category = "task"
)

func (l *TaskList) add(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(l.out, "Missing parameters for \"add\" command.")
		return
	}
	projectName := args[1]
	switch Category(args[0]) {
	case CategoryProject:
		l.addProject(projectName)
	case CategoryTask:
		description := strings.Join(args[2:], " ")
		l.addTask(projectName, Description(description))
	}
}

func (l *TaskList) addProject(name string) {
	l.projectTasks[name] = make([]*Task, 0)
}

func (l *TaskList) addTask(projectName string, description Description) TaskId {
	tasks, ok := l.projectTasks[projectName]
	if !ok {
		fmt.Fprintf(l.out, "Could not find a project with the name \"%s\".\n", projectName)
		return 0
	}
	id := l.nextID()
	l.projectTasks[projectName] = append(tasks, NewTask(id, description, false))
	return id
}

func (l *TaskList) check(idString string) {
	l.setDone(idString, true)
}

func (l *TaskList) uncheck(idString string) {
	l.setDone(idString, false)
}

func (l *TaskList) setDone(idString string, done bool) {
	id, err := newTaskIdFromString(idString)
	if err != nil {
		fmt.Fprintf(l.out, "%s", err)
		return
	}

	for _, tasks := range l.projectTasks {
		for _, task := range tasks {
			if task.GetID() == id {
				task.SetDone(done)
				return
			}
		}
	}

	fmt.Fprintf(l.out, "Task with ID \"%d\" not found.\n", id)
}

func (l *TaskList) nextID() TaskId {
	l.lastID++
	return l.lastID
}

type Deadline struct{ time.Time }

func NewDeadline(y, m, d int) Deadline {
	return Deadline{time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)}
}

func NewDeadlineFromString(date string) (Deadline, error) {
	deadline, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return Deadline{}, err
	}
	return Deadline{deadline}, nil
}

func (l *TaskList) deadline(taskId TaskId, deadline Deadline) {
	for _, tasks := range l.projectTasks {
		for _, task := range tasks {
			if task.id == taskId {
				task.deadline = deadline
			}
		}
	}
}

func (l *TaskList) today() {
	// show projects sequentially
	l.printSortedProjects(l.projectTasks.SortedProjects(), makeFilter(l.timeGetter()))
}

func (l *TaskList) printSortedProjects(sortedProjects []string, filterer filterer) {
	for _, project := range sortedProjects {
		tasks := l.projectTasks[project]
		fmt.Fprintf(l.out, "%s\n", project)
		var filtered []*Task
		for _, task := range tasks {
			if filterer(task) {
				filtered = append(filtered, task)
			}
		}
		for _, task := range filtered {
			done := ' '
			if task.IsDone() {
				done = 'X'
			}
			fmt.Fprintf(l.out, "    [%c] %d: %s\n", done, task.GetID(), task.GetDescription())
		}
		fmt.Fprintln(l.out)
	}
}

type filterer func(*Task) bool

func makeFilter(date time.Time) filterer {
	return func(task *Task) bool {
		return task.deadline.Equal(date)
	}
}

func anyTask(*Task) bool {
	return true
}
