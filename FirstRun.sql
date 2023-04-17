CREATE DATABASE "GolangChat";

\c GolangChat

CREATE TABLE IF NOT EXISTS users
(
    username text NOT NULL,
    password text NOT NULL,
    session_pass text,
    session_time timestamp with time zone,
    id SERIAL NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_key UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS chats
(
    name text NOT NULL,
    id SERIAL NOT NULL,
    CONSTRAINT id PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS chat_users
(
    chat_id integer NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT chat_id FOREIGN KEY (chat_id)
        REFERENCES chats (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages
(
    "time" timestamp with time zone NOT NULL,
    content text NOT NULL,
    chat_id integer NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

CREATE OR REPLACE FUNCTION notify_new_message()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
BEGIN
    PERFORM pg_notify('MessageAdded', json_build_object(
            'chatid', NEW.chat_id,
            'userid', NEW.user_id,
            'content', NEW.content,
            'time', NEW.time,
            'username', users.username
        )::text) FROM messages LEFT JOIN users ON NEW.user_id = users.id;
    RETURN NEW;
END;
$BODY$;

CREATE TRIGGER messageadded
    AFTER INSERT
    ON messages
    FOR EACH ROW
EXECUTE FUNCTION notify_new_message();