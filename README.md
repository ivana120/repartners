# Pack Calculator

A Go API that calculates the optimal number of packs to ship for customer orders.

## Live Demo

http://pack-calculator-alb-1950116497.eu-central-1.elb.amazonaws.com

## Problem

Given configurable pack sizes (e.g., 250, 500, 1000, 2000, 5000), calculate the best combination of packs following these rules:

1. Only whole packs can be sent
2. Minimize total items shipped
3. Minimize number of packs (secondary priority)

### Examples

| Order | Packs | Total |
|-------|-------|-------|
| 1 | 1×250 | 250 |
| 251 | 1×500 | 500 |
| 501 | 1×500 + 1×250 | 750 |
| 12001 | 2×5000 + 1×2000 + 1×250 | 12250 |

### Edge Case

Pack sizes: 23, 31, 53
Order: 500,000
Result: `{23: 2, 31: 7, 53: 9429}` = 500,000 items

## Running Locally

```bash
# With Docker
docker-compose up --build

# Or with Go
make build && make run
```

Access at http://localhost:8080

## API

### Calculate Packs
```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"order_qty": 501}'
```

### Get Pack Sizes
```bash
curl http://localhost:8080/api/pack-sizes
```

### Update Pack Sizes
```bash
curl -X PUT http://localhost:8080/api/pack-sizes \
  -H "Content-Type: application/json" \
  -d '{"pack_sizes": [23, 31, 53]}'
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `PACK_SIZES` | 250,500,1000,2000,5000 | Pack sizes |

## Testing

```bash
make test
```

## Project Structure

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/             # HTTP handlers
│   └── calculator/      # Core algorithm
├── .github/workflows/   # CI/CD pipeline
├── Dockerfile
└── docker-compose.yml
```

## Deployment

The project uses GitHub Actions to automatically deploy to AWS ECS on every push to `main`.

Required GitHub secrets:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
