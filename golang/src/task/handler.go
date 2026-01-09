package main

import (
	"fmt"
	"io"
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

func (h *Handler) add(args []string) {
	if len(args) < 2 {
		_, _ = fmt.Fprintln(h.out, "Missing parameters for \"add\" command.")
		return
	}
	if err := h.TaskList.add(args); err != nil {
		_, _ = fmt.Fprint(h.out, err.Error())
	}
}
