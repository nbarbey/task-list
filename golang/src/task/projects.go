package main

import (
	"fmt"
	"io"
	"sort"
)

type Project []*Task

type Projects map[string]Project

func (p Projects) SortedProjects() []string {
	// sort projects (to make output deterministic)
	sortedProjects := make([]string, 0, len(p))
	for project := range p {
		sortedProjects = append(sortedProjects, project)
	}
	sort.Strings(sortedProjects)

	return sortedProjects
}

func (p Project) Print(w io.Writer) {
	for _, task := range p {
		done := ' '
		if task.IsDone() {
			done = 'X'
		}
		fmt.Fprintf(w, "    [%c] %d: %s\n", done, task.GetID(), task.GetDescription())
	}
	fmt.Fprintln(w)
}

func (p Project) Filter(filterer filtererFunc) Project {
	var filtered []*Task
	for _, task := range p {
		if filterer(task) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}
