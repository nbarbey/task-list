package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskList_today(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	taskList := NewTaskList(in, out)

	taskList.WithCalendar(func() time.Time {
		return time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC)
	})

	taskList.addProject("secrets")
	taskID := taskList.addTask("secrets", "Eat more donuts.")
	taskList.deadline(taskID, NewDeadline(2024, 6, 7))
	taskList.today()

	assert.Equal(t, "secrets\n    [ ] 1: Eat more donuts.\n\n", out.String())
}
