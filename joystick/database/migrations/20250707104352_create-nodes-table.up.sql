create table nodes (
    id serial primary key,
    name text not null,
    ip text not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp
);

alter table service_deployments add column node_id int;
alter table service_deployments add foreign key (node_id) references nodes(id) on delete cascade;
