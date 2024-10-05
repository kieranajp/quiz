CREATE TABLE question_types (
    question_type_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    type_name TEXT -- e.g., 'multiple_choice', 'true_false', etc.
);

CREATE TABLE question_topics (
     topic_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
     topic_name TEXT -- e.g., 'science', 'geography', etc.
);

CREATE TABLE questions (
   question_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
   question_text TEXT,
   question_type_id UUID REFERENCES question_types(question_type_id),
   topic_id UUID REFERENCES question_topics(topic_id),
   trivia TEXT
);

CREATE TABLE answers (
 answer_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
 question_id UUID REFERENCES questions(question_id),
 answer_text TEXT,
 points INT DEFAULT 0 -- Set points for each answer, whether right or wrong
);
