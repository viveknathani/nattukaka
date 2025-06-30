create table if not exists users (
  id serial primary key,
  uuid uuid default gen_random_uuid() not null unique,
  username varchar(255) not null unique,
  password bytea not null
)
