CREATE SCHEMA IF NOT EXISTS chat;

SET SEARCH_PATH TO chat;

CREATE FUNCTION updated_at_trigger() RETURNS TRIGGER
   LANGUAGE plpgsql AS
$$BEGIN
   NEW.updated_at := current_timestamp;
   RETURN NEW;
END;$$;

CREATE TABLE IF NOT EXISTS users (
  id          int           NOT NULL,
  name        VARCHAR(40)   NOT NULL,
  password    VARCHAR(100)  NOT NULL, 
  created_at  timestamptz   NOT NULL  DEFAULT NOW(),
  updated_at  timestamptz   NOT NULL  DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER updated_at_users_trgr
  BEFORE UPDATE
  ON users 
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS conversation (
  id          int           NOT NULL,
  title       VARCHAR(100)   NOT NULL,
  creator_id  int           NOT NULL,
  created_at  timestamptz   NOT NULL  DEFAULT NOW(),
  updated_at  timestamptz   NOT NULL  DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (creator_id) REFERENCES users (id)
);

CREATE TRIGGER updated_at_conversation_trgr
  BEFORE UPDATE
  ON conversation 
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS participants (
  id              int   NOT NULL,
  users_id        int   NOT NULL,
  conversation_id int   NOT NULL,
  created_at      timestamptz   NOT NULL  DEFAULT NOW(),
  updated_at      timestamptz   NOT NULL  DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (users_id) REFERENCES users (id),
  FOREIGN KEY (conversation_id) REFERENCES conversation (id)
);

CREATE TRIGGER updated_at_participants_trgr
  BEFORE UPDATE
  ON participants 
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE IF NOT EXISTS messages (
  guid              VARCHAR(100)  NOT NULL,
  sender_id         int           NOT NULL,
  conversation_id   int           NOT NULL,
  message           VARCHAR(255)  NOT NULL,
  created_at        timestamptz   NOT NULL  DEFAULT NOW(),
  PRIMARY KEY (guid),
  FOREIGN KEY (sender_id) REFERENCES users (id),
  FOREIGN KEY (conversation_id) REFERENCES conversation (id)
);
