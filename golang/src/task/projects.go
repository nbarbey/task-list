package main

import "sort"

type Projects map[string][]*Task

func (p Projects) SortedProjects() []string {
	// sort projects (to make output deterministic)
	sortedProjects := make([]string, 0, len(p))
	for project := range p {
		sortedProjects = append(sortedProjects, project)
	}
	sort.Strings(sortedProjects)

	return sortedProjects
}
