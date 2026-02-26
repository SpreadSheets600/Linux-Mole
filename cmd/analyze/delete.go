package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"

	tea "github.com/charmbracelet/bubbletea"
)

func deletePathCmd(path string, counter *int64) tea.Cmd {
	return func() tea.Msg {
		count, err := trashPathWithProgress(path, counter)
		return deleteProgressMsg{
			done:  true,
			err:   err,
			count: count,
			path:  path,
		}
	}
}

// deleteMultiplePathsCmd moves paths to Trash and aggregates results.
func deleteMultiplePathsCmd(paths []string, counter *int64) tea.Cmd {
	return func() tea.Msg {
		var totalCount int64
		var errors []string

		// Process deeper paths first to avoid parent/child conflicts.
		pathsToDelete := append([]string(nil), paths...)
		sort.Slice(pathsToDelete, func(i, j int) bool {
			return strings.Count(pathsToDelete[i], string(filepath.Separator)) > strings.Count(pathsToDelete[j], string(filepath.Separator))
		})

		for _, path := range pathsToDelete {
			count, err := trashPathWithProgress(path, counter)
			totalCount += count
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				errors = append(errors, err.Error())
			}
		}

		var resultErr error
		if len(errors) > 0 {
			resultErr = &multiDeleteError{errors: errors}
		}

		return deleteProgressMsg{
			done:  true,
			err:   resultErr,
			count: totalCount,
			path:  "",
		}
	}
}

// multiDeleteError holds multiple deletion errors.
type multiDeleteError struct {
	errors []string
}

func (e *multiDeleteError) Error() string {
	if len(e.errors) == 1 {
		return e.errors[0]
	}
	return strings.Join(e.errors[:min(3, len(e.errors))], "; ")
}

func trashPathWithProgress(root string, counter *int64) (int64, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return 0, err
	}

	var count int64
	if info.IsDir() {
		_ = filepath.WalkDir(root, func(_ string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() {
				count++
				if counter != nil {
					atomic.StoreInt64(counter, count)
				}
			}
			return nil
		})
	} else {
		count = 1
		if counter != nil {
			atomic.StoreInt64(counter, 1)
		}
	}

	if err := moveToTrash(root); err != nil {
		return 0, err
	}

	return count, nil
}

func moveToTrash(path string) error {
	return moveToTrashPlatform(path)
}
