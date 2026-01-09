package main

import (
	"fmt"
	"io"
)

type Project []*Task

func (p Project) Print(w io.Writer) {
	for _, task := range p {
		task.Print(w)
	}
	fmt.Fprintln(w)
}

func (p Project) Filter(filterer taskFilterer) Project {
	var filtered []*Task
	for _, task := range p {
		if filterer(task) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}
