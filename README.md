# Bioly Monorepo Setup

This repository contains all **Bioly microservices** and the **Nginx gateway** used to route and secure them.
It is designed as a **monorepo**, where each service lives in its own directory and can be developed, built, and deployed independently — but run together via Docker Compose.

---

## Overview

### Repository Structure

```
.
├── api/                      # OpenAPI documentation
├── common/                   # Common (shared) Go libs
├── configs/                  # Services configs
├── sql/                      # Initial sql scripts for database
├── services/
│   ├── editservice/
│   ├── profileservice/
│   └── authservice/          # Authentication service (Go)                 
├── nginx/                    # Gateway and reverse proxy
│   ├── nginx.conf
│   ├── upstreams.d/
│   ├── locations.d/
│   ├── snippets/
│   └── certs/
├── docker-compose.yml        # Defines all services
├── generate-certs.sh         # Generates local HTTPS certificates
└── README.md                 # This file
```

---

## What’s Inside

The **Nginx gateway** acts as the entry point:
- Terminates HTTPS
- Handles CORS and proxy headers
- Routes requests like:
  - `https://bioly.localhost/auth/...` → `auth` service
  - `https://bioly.localhost/profile/...` → `profile` service

Each service exposes its own internal HTTP port and is registered in `nginx/upstreams.d`.

---

## Prerequisites

Make sure you have:
- **Docker** and **Docker Compose**
- **OpenSSL** (or `mkcert` for trusted local certificates)

---

## 1. Generate Local Certificates

From the project root, run:

```bash
./generate-certs.sh
```

This creates a self-signed certificate for:
- `bioly.localhost`
- `localhost`
- `127.0.0.1`

The files will appear in `nginx/certs/`:
```
nginx/certs/fullchain.pem
nginx/certs/privkey.pem
```

> You can now access APIs at **https://bioly.localhost/** without changing Nginx configs.

---

## 2. Run All Services

Start everything (gateway + backend services):

```bash
docker compose up -d
```

Check logs:

```bash
docker compose logs -f nginx
```

All endpoints are now available under HTTPS.

---

## 3. Test 

#### Health check
```bash
curl -k https://bioly.localhost/auth/ping
# → pong
```

#### Auth login
```bash
curl -k -X POST https://bioly.localhost/auth/login   -H "Content-Type: application/json"   -d '{"username":"admin","password":"secret"}'
```

---

**Enjoy building and running the Bioly microservices stack!**
