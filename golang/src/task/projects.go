package main

import (
	"bytes"
	"fmt"
	"sort"
)

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

func (p Projects) String() string {
	var out string
	for _, projectName := range p.SortedProjects() {
		out += fmt.Sprintf("%s\n", projectName)

		var b bytes.Buffer
		project := p[projectName]
		project.Print(&b)
		out += b.String()
	}

	return out
}
