# log-generator

A Docker container that generates logs at a configurable rate. Useful for load-testing log pipelines (e.g. Loki).

Every log line is exactly 256 bytes.

## Usage

```bash
docker build -t log-generator .
docker run --rm log-generator [rate_KiB_per_sec] [structured|random]
```

### Arguments

| Arg | Default | Description |
|-----|---------|-------------|
| `rate_KiB_per_sec` | `5` | Target output rate in KiB/s |
| `mode` | `structured` | `structured` for JSON logs, `random` for random hex strings |

### Examples

```bash
# 5 KiB/s structured JSON logs (default)
docker run --rm log-generator

# 50 KiB/s structured JSON logs
docker run --rm log-generator 50

# 10 KiB/s random hex strings (harder to compress)
docker run --rm log-generator 10 random
```

### Structured mode

```json
{"ts":"2026-03-17T15:12:13.904Z","level":"DEBUG","i":1,"host":"gen-04","rid":"c1e355a84fa6599f","latency_ms":2251,"bytes_in":13059,"path":"/ready","pad":"xxx..."}
```

### Random mode

```
06e2a039003debdc1a13af4b41a17973ec74666798090dcd687a034f81469ac3cef7383a...
```
