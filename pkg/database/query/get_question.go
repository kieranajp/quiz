package query

import (
	"github.com/google/uuid"
	"github.com/kieranajp/quiz/pkg/database/entity"

	"github.com/jmoiron/sqlx"
)

// GetQuestions preselects questions for a given round
func GetQuestions(db *sqlx.DB, roundType string, numQuestions int) ([]entity.Question, error) {
	var questions []entity.Question

	const preselectQuery = `
		SELECT DISTINCT ON (q.question_id)
			q.question_id,
			q.question_text,
			t.topic_name
		FROM public.questions q
		INNER JOIN public.question_types qt
			ON q.question_type_id = qt.question_type_id
		INNER JOIN public.question_topics t
			ON q.topic_id = t.topic_id
		WHERE qt.type_name = $1
		ORDER BY RANDOM()
		LIMIT $2
	`

	err := db.Select(&questions, preselectQuery, roundType, numQuestions)
	if err != nil {
		return nil, err
	}

	// Fetch answers for each question
	for i, question := range questions {
		answers, err := getAnswersForQuestion(db, question.QuestionID)
		if err != nil {
			return nil, err
		}
		questions[i].Answers = answers
	}

	return questions, nil
}

// getAnswersForQuestion retrieves answers for a given question ID
func getAnswersForQuestion(db *sqlx.DB, questionID uuid.UUID) ([]entity.Answer, error) {
	const answerQuery = `
		SELECT
			a.answer_id,
			a.answer_text
		FROM public.answers a
		WHERE a.question_id = $1
	`

	var answers []entity.Answer
	err := db.Select(&answers, answerQuery, questionID)
	if err != nil {
		return nil, err
	}

	return answers, nil
}
