package query

import (
	"github.com/kieranajp/quiz/pkg/database/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const questionQuery = `
 	SELECT
		q.question_id,
		q.question_text,
		t.topic_name,
		a.answer_id,
		a.answer_text
  	FROM public.questions q
  	INNER JOIN public.answers a
		ON q.question_id = a.question_id
	INNER JOIN public.question_types qt
		ON q.question_type_id = qt.question_type_id
  	INNER JOIN public.question_topics t
		ON q.topic_id = t.topic_id
	WHERE qt.type_name = $1
`

func GetQuestion(db *sqlx.DB, roundType string) (entity.Question, error) {
	var question entity.Question
	answerMap := make(map[uuid.UUID][]entity.Answer)

	type row struct {
		QuestionID   uuid.UUID `db:"question_id"`
		QuestionText string    `db:"question_text"`
		TopicName    string    `db:"topic_name"`
		AnswerID     uuid.UUID `db:"answer_id"`
		AnswerText   string    `db:"answer_text"`
	}

	var rowsData []row
	err := db.Select(&rowsData, questionQuery, roundType)
	if err != nil {
		return entity.Question{}, err
	}

	for _, row := range rowsData {
		if question.QuestionID == uuid.Nil {
			question = entity.Question{
				QuestionID: row.QuestionID,
				Text:       row.QuestionText,
				Topic:      row.TopicName,
				Answers:    []entity.Answer{},
			}
		}

		answerMap[row.QuestionID] = append(answerMap[row.QuestionID], entity.Answer{AnswerID: row.AnswerID, Text: row.AnswerText})
	}

	question.Answers = answerMap[question.QuestionID]

	return question, nil
}
