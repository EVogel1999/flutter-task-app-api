package schema

import "time"

type Task struct {
	ID          string
	Name        string
	Date        time.Time
	Description string
	Category    string
}
