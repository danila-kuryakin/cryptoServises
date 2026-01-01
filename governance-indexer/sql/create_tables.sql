CREATE TABLE proposal (
    id serial primary key,
    hex_id varchar(300) unique not null ,
    created_at timestamp,
    state text,
    space_id serial,
    Author varchar(300)
);

CREATE TABLE space (
    id serial primary key,
    hex_id varchar(300),
    name varchar(300)
);