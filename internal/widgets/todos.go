package widgets

import (
	"fmt"

	"github.com/namyoungkim/visor/internal/config"
	"github.com/namyoungkim/visor/internal/input"
	"github.com/namyoungkim/visor/internal/render"
	"github.com/namyoungkim/visor/internal/transcript"
)

// TodosWidget displays task progress from TaskCreate/TaskUpdate tools.
//
// Supported Extra options:
//   - show_label: "true"/"false" - show "Tasks:" prefix (default: false)
//   - max_subject_len: maximum length for task subject (default: 30)
//
// Output format:
//   - "✓ All done (5/5)" when all tasks completed
//   - "⊙ Task name (3/5)" when tasks in progress
//   - "○ Task name (0/5)" when all tasks pending
type TodosWidget struct {
	transcript *transcript.Data
}

func (w *TodosWidget) Name() string {
	return "todos"
}

// SetTranscript sets the transcript data for this widget.
func (w *TodosWidget) SetTranscript(t *transcript.Data) {
	w.transcript = t
}

func (w *TodosWidget) Render(session *input.Session, cfg *config.WidgetConfig) string {
	if w.transcript == nil || len(w.transcript.Todos) == 0 {
		return ""
	}

	todos := w.transcript.Todos
	total := len(todos)

	// Count by status
	completed := 0
	inProgress := 0
	var currentTask *transcript.Todo

	for i := range todos {
		switch todos[i].Status {
		case transcript.TodoCompleted:
			completed++
		case transcript.TodoInProgress:
			inProgress++
			if currentTask == nil {
				currentTask = &todos[i]
			}
		case transcript.TodoPending:
			if currentTask == nil && inProgress == 0 {
				currentTask = &todos[i]
			}
		}
	}

	maxSubjectLen := GetExtraInt(cfg, "max_subject_len", 30)
	var text string
	var color string

	if completed == total {
		// All done
		text = fmt.Sprintf("✓ All done (%d/%d)", completed, total)
		color = "green"
	} else if currentTask != nil {
		// Show current task
		subject := truncateString(currentTask.Subject, maxSubjectLen)
		icon := "○" // pending
		color = "yellow"
		if currentTask.Status == transcript.TodoInProgress {
			icon = "⊙"
			color = "cyan"
		}
		text = fmt.Sprintf("%s %s (%d/%d)", icon, subject, completed, total)
	} else {
		// Fallback
		text = fmt.Sprintf("○ Tasks (%d/%d)", completed, total)
		color = "yellow"
	}

	if GetExtraBool(cfg, "show_label", false) {
		text = "Tasks: " + text
	}

	return render.Colorize(text, color)
}

func (w *TodosWidget) ShouldRender(session *input.Session, cfg *config.WidgetConfig) bool {
	return w.transcript != nil && len(w.transcript.Todos) > 0
}
