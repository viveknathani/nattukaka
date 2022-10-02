-- "name" is stored as varchar to support variable length string with 
-- no upper limit. "password" is stored as bytea to support the bcrypt implementation. 
-- The Go implementation outputs the hash in base64 format which makes varchar an alternative
-- to bytea in our case. I still choose to go with bytea to be verbose about the fact that we
-- are dealing with binary data here. "email" is stored in accordance with RFC 5321 and Errata
-- 1690.
create table if not exists users(
    id uuid primary key,
    name varchar not null,
    email varchar(319) not null,
    password bytea not null
);

create type taskStatus as enum('pending', 'done');

create table if not exists todos(
    id uuid primary key,
    userId uuid references users(id),
    task varchar not null,
    status taskStatus,
    deadline date,
    completedAt date
);

create table if not exists notes(
    id uuid primary key,
    userId uuid references users(id),
    title varchar not null,
    content text
);

create table if not exists logs(
    level varchar,
    ts real,
    path varchar,
    message text,
    requestID varchar,
    method varchar,
    ip varchar,
    info jsonb
);
