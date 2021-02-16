package client

type Message struct {
	TypeMessage string
	QuestionId  int
	AnswerIds   []int
}

type Question struct {
	QuestionId  int
	AnswerIds   []int
}

type Answer struct {
	QuestionId int `json:"questionId,omitempty"`
	AnswerId   int `json:"answerId,omitempty"`
}
