# nattukaka - design

At it's heart, nattukaka is supposed to be a container provisioning system.

We'd run it on a single node to save costs and stay lean but the design should theoretically work on a multi-node setup as well.

We come up with two services:
1. joystick - interacts with the DB + gives a REST API + speaks to player service when needed
2. player - does container orchestration and manages the lifecycle of containers + gives a GRPC API

Theoretically, there can be multiple instances of the player service across multiple nodes.

### db schemas

`users` table
- id
- username
- password

`services` table
- id
- uuid
- name
- repository_url
- branch
- env_vars
- owner_id
- created_at
- updated_at

`service_deployments` table
- id
- uuid
- service_id
- commit
- status
- container_id
- created_at

### API (joystick)
- `POST /services`
- `GET /services`
- `GET /services/:id`
- `PUT /services/:id`
- `DELETE /services/:id`
- `GET /services/:id/deployments`
- `POST /services/:id/deployments`
- `GET /services/:id/deployments/:id`
- `PATCH /services/:id/deployments/:id` (only player service can call this)

### API (player)
- buildImage
- startContainer
- stopContainer
- streamLogs

### service deployment states

- QUEUED
- BUILDING
- STARTING
- RUNNING
- STOPPING
- STOPPED
- FAILED

### port allocation

Done by the player service - randomly picking ports from a range.

### caddy updates

Caddy will be used to route public traffic to the running container

Route format: `${service_name}.nattukaka.dev → localhost:<allocated_port>`
