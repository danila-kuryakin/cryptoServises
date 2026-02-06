CREATE TABLE IF NOT EXISTS proposals
(
    id serial primary key,
    hex_id varchar(256) unique not null ,
    title text,
    author varchar(256),
    created_at integer,
    start_at integer,
    end_at integer,
    snapshot bigint,
    state text,
    choices json,
    space_id varchar(256),
    space_name varchar(256)
);

CREATE TABLE IF NOT EXISTS proposals_outbox
(
    id serial primary key,
    hex_id varchar(256) unique not null references proposals(hex_id) on delete cascade,
    event_type text,
    created_at timestamp,
    processed_at timestamp NULL
);

CREATE TABLE IF NOT EXISTS event_scheduler
(
    id serial primary key,
    hex_id varchar(256) not null references proposals(hex_id) on delete cascade,
    event_type varchar(256),
    event_at timestamp not null,
    processed_at timestamp NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id serial primary key,
    user_id bigint unique not null,
    username text,
    dao_subscribed integer
);

CREATE TABLE IF NOT EXISTS event_outbox
(
    id serial primary key,
    hex_id varchar(256) unique not null references proposals(hex_id) on delete cascade,
    event_type text,
    created_at timestamp,
    processed_at timestamp NULL
);