CREATE SCHEMA IF NOT EXISTS chat;

SET SEARCH_PATH TO chat;

CREATE FUNCTION updated_at_trigger() RETURNS TRIGGER
   LANGUAGE plpgsql AS
$$BEGIN
   NEW.updated_at := current_timestamp;
   RETURN NEW;
END;$$;

CREATE TABLE IF NOT EXISTS users (
  id                bigserial     NOT NULL,
  name              VARCHAR(40)   NOT NULL,
  login             VARCHAR(40)   NOT NULL,
  color             VARCHAR(7)    NOT NULL,
  password_hash     VARCHAR(60)  NOT NULL, 
  created_at        timestamptz   NOT NULL  DEFAULT NOW(),
  updated_at        timestamptz   NOT NULL  DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER updated_at_users_trgr
  BEFORE UPDATE
  ON users 
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS conversations (
  id                  bigserial             NOT NULL,
  title               VARCHAR(100)          NOT NULL,
  conversation_kind   int                   NOT NULL,
  creator_id          bigserial             NOT NULL,
  created_at          timestamptz           NOT NULL  DEFAULT NOW(),
  updated_at          timestamptz           NOT NULL  DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (creator_id) REFERENCES users (id)
);

CREATE TRIGGER updated_at_conversations_trgr
  BEFORE UPDATE
  ON conversations 
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS participants (
  id              bigserial     NOT NULL,
  user_id         bigserial     NOT NULL,
  conversation_id bigserial     NOT NULL,
  created_at      timestamptz   NOT NULL  DEFAULT NOW(),
  updated_at      timestamptz   NOT NULL  DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (conversation_id) REFERENCES conversations (id)
);

CREATE TRIGGER updated_at_participants_trgr
  BEFORE UPDATE
  ON participants 
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS messages (
  id                uuid            NOT NULL,
  sender_id         bigserial       NOT NULL,
  conversation_id   bigserial       NOT NULL,
  message_kind      int             NOT NULL,
  message           VARCHAR(1024)   NOT NULL,
  created_at        timestamptz     NOT NULL  DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (sender_id) REFERENCES users (id),
  FOREIGN KEY (conversation_id) REFERENCES conversations (id)
);

CREATE TABLE IF NOT EXISTS friends (
  id          bigserial   NOT NULL,
  user_id     bigserial   NOT NULL,
  friend_id   bigserial   NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (friend_id) REFERENCES users (id)
);

CREATE TRIGGER updated_at_friends_trgr
  BEFORE UPDATE
  ON friends 
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();
