create table if not exists services (
  id serial primary key,
  uuid uuid default gen_random_uuid() not null unique,
  name varchar(255) not null unique,
  repository_url varchar(255) not null,
  branch varchar(255) not null,
  env_vars jsonb not null,
  port_mapping jsonb not null,
  owner_id int not null,
  created_at timestamp not null,
  updated_at timestamp not null
);