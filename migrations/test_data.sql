SET SEARCH_PATH TO chat;

INSERT INTO users (id, name, password) VALUES (1, 'Alice', 'password123');
INSERT INTO users (id, name, password) VALUES (2, 'Bob', 'password456');
INSERT INTO users (id, name, password) VALUES (3, 'Charlie', 'password789');

INSERT INTO conversation (id, title, creator_id) VALUES (1, 'Project Discussion', 1);
INSERT INTO conversation (id, title, creator_id) VALUES (2, 'General Chat', 2);

INSERT INTO messages (guid, sender_id, conversation_id, message) VALUES ('a1b2c3', 1, 1, 'Hello everyone, let''s start the project discussion.');
INSERT INTO messages (guid, sender_id, conversation_id, message) VALUES ('d4e5f6', 2, 1, 'Sure, I am ready.');
INSERT INTO messages (guid, sender_id, conversation_id, message) VALUES ('g7h8i9', 1, 2, 'How was your day, Bob?');
INSERT INTO messages (guid, sender_id, conversation_id, message) VALUES ('j0k1l2', 2, 2, 'Pretty good, thanks for asking.');

-- Alice and Bob in 'Project Discussion'
INSERT INTO participants (id, users_id, conversation_id) VALUES (1, 1, 1);
INSERT INTO participants (id, users_id, conversation_id) VALUES (2, 2, 1);
-- Alice and Bob in 'General Chat'
INSERT INTO participants (id, users_id, conversation_id) VALUES (3, 1, 2);
INSERT INTO participants (id, users_id, conversation_id) VALUES (4, 2, 2);
-- Adding Charlie to 'General Chat'
INSERT INTO participants (id, users_id, conversation_id) VALUES (5, 3, 2);
