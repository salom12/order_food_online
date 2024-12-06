# Order Food Online API

This is a Go-based API for managing products, orders, and promo codes. It uses PostgreSQL as the primary database and Redis for caching.

## Features
- **Products**: Fetch a list of products.
- **Orders**: Place orders with optional promo codes.
- **Promo Codes**: Validate promo codes using predefined rules.


## Getting Started
### Prerequisites
- Docker
- Go 1.20+

### Run Locally
```bash
docker-compose up --build
