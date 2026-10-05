package pomodoro

import "time"

const (
	studyBlock = 25
	breakBlock = 5
	tick       = time.Minute
)

type segment struct {
	label   string
	minutes int
}

func buildPlan(total int) []segment {
	var plan []segment
	remaining := total

	for remaining > 0 {
		study := min(studyBlock, remaining)
		remaining -= study

		plan = append(plan, segment{"STUDY", study})

		if remaining > 0 {
			plan = append(plan, segment{"BREAK", breakBlock})
		}
	}

	return plan
}
