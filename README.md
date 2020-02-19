## Machinable Event Processor

This project processes and sends configured Machinable project web hooks via HTTP POST.

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
