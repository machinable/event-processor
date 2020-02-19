## Machinable Event Processor

Processes and sends configured Machinable project web hooks via HTTP POST.

Web Hooks are sent as the following sample payload:

```json
{
    "entity_key": "cars",
    "entity_type": "resource",
    "event": "create",
    "payload": {
        "_metadata": {
            "created": 1582072623,
            "creator": "anonymous",
            "creator_type": "anonymous"
        },
        "id": "2a5cf445-ccdd-4d55-a6d6-2fa25ab2ba58",
        "make": "Chevy",
        "model": "Farquaid",
        "year": 2020
    },
    "project_Id": "5c0f744b-65ee-4c62-b257-55dd546c7f1e"
}
```

|Field|Description|
|-----|-----------|
|`entity_key`|The name of the API Resource or JSON Key|
|`entity_type`|`resource` or `json`|
|`event`|The type of event that triggered the Web Hook, i.e. `create`, `edit`, or `delete`|
|`payload`|The JSON payload of the request that triggered the Web Hook|
|`project_id`|The ID of the Machinable project|


The result of the webhook HTTP request is then placed on another Redis queue for processing by the Machinable API.

### Run Locally

```sh
# build without cache
rebuild:
	docker-compose build --no-cache

# build image
build:
	docker-compose build

# spin up containers (postgres and api)
up:
	docker-compose up -d

# bring down running containers
down:
	docker-compose down

# stop running containers
stop:
	docker-compose stop

# remove containers and images
remove:
	docker-compose rm -f

# cleanup images and volumes
clean:
	docker-compose down --rmi all -v --remove-orphans
```
