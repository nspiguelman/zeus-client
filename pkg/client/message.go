package client

type Message struct {
	TypeMessage string
}

type Question struct {
	TypeMessage string
	QuestionId  int
	AnswerIds   []int
}

type Answer struct {
	QuestionId int `json:"questionId,omitempty"`
	AnswerId   int `json:"answerId,omitempty"`
}
