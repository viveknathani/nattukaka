-- users table
create table if not exists users (
    id serial primary key,
    public_id varchar(255) unique not null,
    name varchar(255) not null,
    email varchar(255) not null
);

-- workspaces table
create table if not exists workspaces (
    id serial primary key,
    public_id varchar(255) unique not null,
    name varchar(255) not null
);

-- workspace_users table
create table if not exists workspace_users (
    id serial primary key,
    public_id varchar(255) unique not null,
    workspace_id integer not null references workspaces(id) on delete cascade,
    user_id integer not null references users(id) on delete cascade,
    role varchar(255) not null
);

-- instance_types table
create table if not exists instance_types (
    id serial primary key,
    public_id varchar(255) unique not null,
    name varchar(255) not null,
    cpu real not null,
    memory real not null,
    disk real not null
);

-- services table
create table if not exists services (
    id serial primary key,
    public_id varchar(255) unique not null,
    name varchar(255) not null,
    status varchar(255) not null,
    type varchar(255) not null,
    runtime varchar(255),
    workspace_id integer not null references workspaces(id) on delete cascade,
    created_by integer not null references users(id) on delete cascade,
    last_deployed_at timestamp,
    created_at timestamp default current_timestamp,
    instance_type_id integer not null references instance_types(id),
    internal_url varchar(2048),
    external_url varchar(2048)
);

-- web_services table
create table if not exists web_services (
    id serial primary key,
    public_id varchar(255) unique not null,
    service_id integer not null references services(id) on delete cascade,
    repository varchar(255),
    branch varchar(255),
    root_directory varchar(255),
    build_command text,
    pre_deploy_command text,
    start_command text,
    health_check_path varchar(255),
    environment text
);

-- database_services table
create table if not exists database_services (
    id serial primary key,
    public_id varchar(255) unique not null,
    service_id integer not null references services(id) on delete cascade
);

-- volumes table
create table if not exists volumes (
    id serial primary key,
    public_id varchar(255) unique not null
);

-- service_volumes table
create table if not exists service_volumes (
    service_id integer not null references services(id) on delete cascade,
    volume_id integer not null references volumes(id) on delete cascade,
    primary key (service_id, volume_id)
);

-- deploys table
create table if not exists deploys (
    id serial primary key,
    public_id varchar(255) unique not null,
    service_id integer not null references services(id) on delete cascade,
    status varchar(50) not null,
    commit varchar(255),
    image varchar(255)
);
