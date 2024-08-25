package models

type Result struct {
	ID         string `json:"id"`
	RunnerID   string `json:"runner_id"`
	RaceResutl string `json:"race_result"`
	Location   string `json:"location"`
	Position   int    `json:"position"`
	Year       int    `json:"year"`
}
