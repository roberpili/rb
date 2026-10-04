package pomodoro

import "time"

const (
	studyBlock = 25
	breakBlock = 5
	tick       = time.Second
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

		// extend if its not enough for another round
		if remaining <= breakBlock {
			study += remaining
			remaining = 0
		}
		plan = append(plan, segment{"STUDY", study})

		if remaining > 0 {
			plan = append(plan, segment{"BREAK", breakBlock})
			remaining -= breakBlock
		}
	}

	return plan
}
