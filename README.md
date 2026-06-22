# HTTP to MQTT Bridge

A minimal HTTP to MQTT bridge written in Go that allows publishing MQTT messages via HTTP GET or POST requests.

## Environment Variables

| Variable | Required | Default | Description |
| :--- | :---: | :---: | :--- |
| `AUTH_KEY` | **Yes** | — | Secret token required to authorize requests. |
| `MQTT_HOST` | **Yes** | — | Broker connection URI (e.g., `tcp://localhost:1883`). |
| `MQTT_USER` | No | — | Username for broker authentication. |
| `MQTT_PASS` | No | — | Password for broker authentication. |
| `HTTP_PORT` | No | `8080` | Port for the HTTP server to listen on. |

---

## API Reference

### Endpoint
`ANY /publish`

### Parameters
Parameters must be passed as URL query variables (GET) or form-encoded body parameters (POST). Headers are ignored for authentication.

* `auth_key` (Required): Must match the application's `AUTH_KEY` environment variable.
* `topic` (Required): Target MQTT topic.
* `payload` (Required): String content to publish.
* `qos` (Optional): MQTT Quality of Service level (`0`, `1`, or `2`). Defaults to `1`.

---

## Examples

### HTTP GET
```bash
curl "http://localhost:8080/publish?auth_key=secret&topic=home/test&payload=hello&qos=0"
