# Caching Proxy Server

A simple HTTP caching proxy server built in Go that caches responses from origin servers for faster subsequent requests.

## Installation

```bash
git clone https://github.com/AmrMurad1/cache-proxy.git
cd cache-proxy
go build -o caching-proxy
```

## Usage

### Start Server
```bash
./caching-proxy --port 3000 --origin http://dummyjson.com
```

### Make Requests
```bash
# First request - Cache MISS
curl http://localhost:3000/products

# Second request - Cache HIT (faster!)
curl http://localhost:3000/products
```

### Clear Cache
```bash
./caching-proxy --clear-cache
```

## How It Works

Requests are forwarded to the origin server and cached. Subsequent identical requests are served from cache. Response includes `X-Cache: HIT` or `X-Cache: MISS` header.

## Options

- `--port` - Port number for proxy server (required)
- `--origin` - Origin server URL (required)
- `--clear-cache` - Clear all cached responses

## Requirements

Go 1.16+

## Project from

[roadmap.sh - Caching Server Project](https://roadmap.sh/projects/caching-server)
