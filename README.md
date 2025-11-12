# Inkstack

## What is Inkstack?

Inkstack is a full-stack, production-grade web app you’ll build to sharpen engineering skills and showcase real-world competencies.
The first incarnation is a developer blog & knowledge hub with accounts, posts, comments, tags, search, media, and an admin dashboard—designed to be incrementally evolved into a distributed system.

## Purpose

This project serves as a platform to introduce and practice cutting-edge technologies step by step, progressing from simple concepts to advanced implementations, with the goal of mastering these technologies.
- Frontend
  - Next.js
- Backend
  - Gin
  - Authentication: OAuth 2
  - Storage:
    - Relational: PostgreSQL
    - NoSQL: MongoDB, Redis (cache/session)
  - Search: Elasticsearch / OpenSearch
  - Queues & Async Jobs: RabbitMQ, Kafka
- Architecture
  - Microservices
- DevOps
  - Containerization: Docker, K8s
  - CI/CD: GitHub Actions
- Observability & Reliability
  - Monitoring & Metrics: Prometheus + Grafana
  - Logging: ELK stack (Elasticsearch, Logstash, Kibana) or Loki
- Security & Compliance
  - HTTPS/TLS, HSTS, secure cookies
  - Role-Based Access Control (RBAC)


## README-Driven Coding

- Development happens in iterations.
- The project is documented with high-level descriptions, feature introductions, and functional catalogs. These documents give the AI the necessary background knowledge and context. 
- For each iteration, a new TODO markdown is written to specify what the AI should implement or modify. 
- The AI reads both the project documentation and the TODO file, then generates the corresponding code. 
- The developer tests and adjusts the generated code as needed. 
- After adjustments, the AI reviews the modified code and updates the documentation accordingly, ensuring that the next iteration proceeds smoothly.


# Inkstack Microservices Architecture

## Overview

Inkstack is now architected as a microservices system with two main services:

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Client    │────────▶│ Auth Service │────────▶│  Auth DB    │
│             │         │   :8082      │         │   :5433     │
└─────────────┘         └──────────────┘         └─────────────┘
       │                        │
       │                   JWT Token
       │                        │
       ▼                        ▼
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│  Validates  │────────▶│ API Service  │────────▶│  API DB     │
│  JWT Token  │         │   :8081      │         │   :5432     │
└─────────────┘         └──────────────┘         └─────────────┘
                               │
                               ▼
                        ┌─────────────┐
                        │   Redis     │
                        │   :6379     │
                        └─────────────┘
```

## Services

### 1. API Service (Port 8081)
**Responsibility:** Business logic and data management
- Posts CRUD
- Comments CRUD
- Tags, search, media (future)

**Database:** `inkstack_api` (Port 5432)
- Posts table
- Comments table
- Tags, media (future)

**Tech Stack:**
- Go + Gin
- PostgreSQL
- GORM

### 2. Auth Service (Port 8082)
**Responsibility:** User authentication and authorization
- User registration & login
- JWT token generation & validation
- Password management
- OAuth 2.0 integration (future)
- Token refresh & revocation

**Database:** `inkstack_auth` (Port 5433)
- Users table
- Refresh tokens
- OAuth accounts (future)

**Tech Stack:**
- Go + Gin
- PostgreSQL
- Redis (token blacklist, rate limiting)
- JWT (github.com/golang-jwt/jwt/v5)
- Bcrypt (password hashing)

### 3. Redis (Port 6379)
**Responsibility:** Caching and temporary data
- JWT token blacklist
- Rate limiting counters
- Session cache

