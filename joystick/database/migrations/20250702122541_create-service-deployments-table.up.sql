create table if not exists service_deployments (
    id serial primary key,
    uuid uuid default gen_random_uuid() not null unique,
    service_id int not null,
    commit varchar(255) not null,
    status varchar(255) not null,
    container_id varchar(255) not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    foreign key (service_id) references services(id) on delete cascade
);