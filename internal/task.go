package waitter

import (
	"fmt"
	"io"
	"os"
)

// Task represents a task to be executed to a file
type Task struct {
	Order      int                 `json:"order"`
	Type       string              `json:"type"`
	Parameters []map[string]string `json:"parameters"`
}

// TaskTypeFileRemove represents a remove task for a file
type TaskTypeFileRemove struct {
	file string
}

type ByOrder []Task

func (a ByOrder) Len() int           { return len(a) }
func (a ByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }
func (a ByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Remove arquivo
func (r *TaskTypeFileRemove) Remove() error {
	err := os.Remove(r.file)
	if err != nil {
		return err
	}
	return nil
}

// TaskTypeFileCopy Representa uma tarefa de copia
type TaskTypeFileCopy struct {
	source      string
	destination string
}

// Copy files
func (c *TaskTypeFileCopy) Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
