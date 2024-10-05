package entity

import "github.com/google/uuid"

// Question struct to hold question data
type Question struct {
	QuestionID uuid.UUID `json:"question_id"`
	Text       string    `json:"text"`
	Topic      string    `json:"topic"`
	Answers    []Answer  `json:"answers"`
}

// Answer struct to hold answer data
type Answer struct {
	AnswerID uuid.UUID `json:"answer_id"`
	Text     string    `json:"text"`
}
