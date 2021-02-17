package client

type Message struct {
	TypeMessage string
	QuestionId  int
	AnswerIds   []int
	Scores 		map[string]Score
}

type Question struct {
	QuestionId  int
	AnswerIds   []int
}

type Score struct {
	Score     int
	IsCorrect bool
}

type Answer struct {
	QuestionId int `json:"questionId,omitempty"`
	AnswerId   int `json:"answerId,omitempty"`
}
