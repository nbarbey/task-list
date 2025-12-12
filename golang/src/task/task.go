package main

import (
	"fmt"
	"io"
	"strconv"
)

type TaskId int64

func newTaskIdFromString(idString string) (TaskId, error) {
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid ID \"%s\".\n", idString)
	}
	return TaskId(id), nil
}

type Description string

// Task describes an elementary task.
type Task struct {
	id          TaskId
	description Description
	done        bool
	deadline    Deadline
}

// NewTask initializes a Task with the given ID, description and completion status.
func NewTask(id TaskId, description Description, done bool) *Task {
	return &Task{
		id:          id,
		description: Description(description),
		done:        done,
	}
}

// GetID returns the task ID.
func (t *Task) GetID() TaskId {
	return t.id
}

// GetDescription returns the task description.
func (t *Task) GetDescription() Description {
	return t.description
}

// IsDone returns whether the task is done or not.
func (t *Task) IsDone() bool {
	return t.done
}

// Print prints the task as a string to the writer
func (t *Task) Print(w io.Writer) {
	fmt.Fprint(w, t.Sprint())
}

// Spring returns a string representation of the Task
func (t *Task) Sprint() string {
	done := ' '
	if t.IsDone() {
		done = 'X'
	}
	return fmt.Sprintf("    [%c] %d: %s\n", done, t.GetID(), t.GetDescription())
}

// SetDone changes the completion status of the task.
func (t *Task) SetDone(done bool) {
	t.done = done
}
