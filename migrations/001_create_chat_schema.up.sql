CREATE SCHEMA IF NOT EXISTS chat;

SET SEARCH_PATH TO chat;

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

CREATE TABLE IF NOT EXISTS messages (
  id                uuid            NOT NULL,
  sender_id         bigserial       NOT NULL,
  message_kind      int             NOT NULL,
  message           VARCHAR(1024)   NOT NULL,
  created_at        timestamptz     NOT NULL  DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (sender_id) REFERENCES users (id)
);
