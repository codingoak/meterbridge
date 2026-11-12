# meterbridge

A small, **read-only** gateway that polls energy meters over **Modbus TCP** and exposes the readings via a **REST API** (and optionally **MQTT**).

> This is my first Go project. I'm building it to learn Go while working on a problem I know from my background in electrical installation and building automation: getting measurement data out of field devices safely and reliably.

**Status:** early development – see [Roadmap](#roadmap) for what already works.

---

## What it does

```text
┌──────────────┐  Modbus TCP   ┌──────────────┐   REST (JSON)   ┌──────────────┐
│ Energy meter │ ◄──────────── │  meterbridge │ ──────────────► │   Clients    │
│(or simulator)│   read-only   │   (poller)   │                 │ dashboards…  │
└──────────────┘               └──────┬───────┘                 └──────────────┘
                                      │ optional
                                      ▼ MQTT
                               ┌──────────────┐
                               │  Mosquitto   │
                               └──────────────┘
```

- Polls configured registers (voltage, current, power, energy) at a fixed interval
- Applies scaling and units from a YAML config
- Serves the latest values via `GET /api/values` and a health check via `GET /api/health`
- Includes a small **Modbus TCP simulator**, so no hardware is needed for development and tests

## Simulator register map

Holding registers, float32, big endian, high word first.

| Address | Value   | Unit |
| ------- | ------- | ---- |
| 0–1     | Voltage | V    |
| 2–3     | Current | A    |
| 4–5     | Power   | W    |
| 6–7     | Energy  | kWh  |

Reads outside this range return exception code 2 (Illegal Data Address).

## Why read-only

Modbus TCP has **no authentication and no encryption**. Anyone who can reach the device on the network can read – and write – registers. In building automation and other OT environments, an unintended write can change setpoints or switch equipment.

meterbridge therefore:

- **never writes** to any register (only read function codes are used)
- will validate the config at startup (allowed register ranges, data types)
- is meant to run in a **segmented network**, with only the gateway allowed to reach the field devices

The simulator follows the same rule: write requests are rejected with exception code 1 (Illegal Function).

## Getting started

Requires Go 1.26+.

```sh
# Start the simulator (Modbus TCP on port 5502)
go run ./cmd/simulator
```

Read the values with [mbpoll](https://github.com/epsilonrt/mbpoll):

```sh
mbpoll -a 1 -r 1 -c 4 -t 4:float -B -p 5502 -1 127.0.0.1
```

`-B` sets big-endian word order (high word first), matching the register map above.

The gateway (`cmd/meterbridge`) and the REST API are not implemented yet.

## Project structure

```text
cmd/
  meterbridge/ # gateway entry point
  simulator/ # Modbus TCP meter simulator
internal/
  modbus/ # polling, register parsing, scaling
  api/ # REST handlers
config.example.yaml
```

## Roadmap

- [x] Modbus TCP simulator with static meter values
- [ ] Simulator: fluctuating values
- [ ] Poller: read registers from YAML config, interval, timeout and reconnect
- [ ] REST API: `/api/values`, `/api/health` (net/http, no framework)
- [ ] Unit tests (scaling, parsing) and integration test against the simulator
- [ ] CI with GitHub Actions (`go test`, `go vet`)
- [ ] Optional: MQTT publishing and `docker-compose.yml` (simulator, gateway, broker)

## Tech

- Go (standard library first: `net/http`, `context`, `time`)
- Modbus: [simonvetter/modbus](https://github.com/simonvetter/modbus)
- MQTT (planned): [eclipse/paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang)

## License

MIT – see [LICENSE](LICENSE).
