INSERT INTO questions (question_id, question_text, question_type_id, topic_id, trivia) VALUES
                                                                                           (uuid_generate_v4(), 'What is the heaviest noble gas?', 'a8be954f-2ce0-4e82-8f46-a3aad89c2f2f', 'cd631879-a2c7-460b-95c6-f2f371d045c5', 'Radon is the heaviest noble gas and is radioactive. It occurs naturally as a decay product of uranium and can be found in some soils and homes.'),
                                                                                           (uuid_generate_v4(), 'What is the largest internal organ in the human body?', 'a8be954f-2ce0-4e82-8f46-a3aad89c2f2f', 'cd631879-a2c7-460b-95c6-f2f371d045c5', 'The liver is the largest internal organ and has over 500 vital functions, including detoxifying chemicals and metabolizing drugs.'),
                                                                                           (uuid_generate_v4(), 'Which planet has the shortest day?', 'a8be954f-2ce0-4e82-8f46-a3aad89c2f2f', 'cd631879-a2c7-460b-95c6-f2f371d045c5', 'Jupiter has the shortest day of all the planets in our solar system, completing one rotation in just under 10 hours.'),
                                                                                           (uuid_generate_v4(), 'What is the chemical name for laughing gas?', 'a8be954f-2ce0-4e82-8f46-a3aad89c2f2f', 'cd631879-a2c7-460b-95c6-f2f371d045c5', 'Nitrous oxide, commonly known as laughing gas, is used as an anaesthetic in dentistry and surgery due to its pain-relieving properties.'),
                                                                                           (uuid_generate_v4(), 'How many bones are in the adult human body?', 'a8be954f-2ce0-4e82-8f46-a3aad89c2f2f', 'cd631879-a2c7-460b-95c6-f2f371d045c5', 'The adult human body has 206 bones, but babies are born with around 270 bones, many of which fuse together as they grow.');

INSERT INTO answers (answer_id, question_id, answer_text, points) VALUES
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the heaviest noble gas?'), 'Argon', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the heaviest noble gas?'), 'Xenon', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the heaviest noble gas?'), 'Radon', 100),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the heaviest noble gas?'), 'Helium', 0),

                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the largest internal organ in the human body?'), 'Liver', 100),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the largest internal organ in the human body?'), 'Lungs', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the largest internal organ in the human body?'), 'Heart', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the largest internal organ in the human body?'), 'Kidney', 0),

                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'Which planet has the shortest day?'), 'Jupiter', 100),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'Which planet has the shortest day?'), 'Mars', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'Which planet has the shortest day?'), 'Venus', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'Which planet has the shortest day?'), 'Mercury', 0),

                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the chemical name for laughing gas?'), 'Nitrogen dioxide', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the chemical name for laughing gas?'), 'Nitrous oxide', 100),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the chemical name for laughing gas?'), 'Carbon monoxide', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'What is the chemical name for laughing gas?'), 'Sulfur dioxide', 0),

                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'How many bones are in the adult human body?'), '196', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'How many bones are in the adult human body?'), '206', 100),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'How many bones are in the adult human body?'), '220', 0),
                                                                      (uuid_generate_v4(), (SELECT question_id FROM questions WHERE question_text = 'How many bones are in the adult human body?'), '176', 0);
