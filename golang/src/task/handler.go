package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Handler struct {
	in  io.Reader
	out io.Writer
	*TaskList
}

func (h *Handler) WithCalendar(getter TimeGetter) *Handler {
	h.TaskList = h.TaskList.WithCalendar(getter)
	return h
}

func NewHandler(in io.Reader, out io.Writer) *Handler {
	return NewHandlerFor(NewTaskList(in, out), in, out)
}

func NewHandlerFor(list *TaskList, in io.Reader, out io.Writer) *Handler {
	return &Handler{in: in, out: out, TaskList: list}
}

// Run runs the command loop of the task manager.
// Sequentially executes any given command, until the user types the Quit message.
func (l *Handler) Run() {
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

func (l *Handler) execute(cmdLine string) {
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

func (h *Handler) add(args []string) {
	if len(args) < 2 {
		_, _ = fmt.Fprintln(h.out, "Missing parameters for \"add\" command.")
		return
	}
	if err := h.TaskList.add(args); err != nil {
		_, _ = fmt.Fprint(h.out, err.Error())
	}
}

func (h *Handler) today() {
	todayProjects := h.TaskList.today()
	formattedProjects := todayProjects.String()
	fmt.Fprint(h.out, formattedProjects)
}
