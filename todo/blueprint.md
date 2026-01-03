# Inkstack

## What is Inkstack?

Inkstack is a full-stack, production-grade web app you’ll build to sharpen engineering skills and showcase real-world competencies.
The first incarnation is a developer blog & knowledge hub with accounts, posts, comments, tags, search, media, and an admin dashboard—designed to be incrementally evolved into a distributed system.

## Purpose

This project serves as a platform to introduce and practice cutting-edge technologies step by step, progressing from simple concepts to advanced implementations, with the goal of mastering these technologies.
- Frontend
  - Vue.js
- Backend
  - Gin
  - Authentication: OAuth 2
  - Storage:
    - Relational: PostgreSQL
    - NoSQL: MongoDB, Redis (cache/session)
  - Swagger
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
