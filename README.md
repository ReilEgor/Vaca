# Vaca
  Vaca is the primary repository for a microservices-based job aggregator platform specifically tailored for the Ukrainian tech market.
It represents a full-cycle engineering solution that automates the collection, processing, and searching of job listings from various local platforms.

**Core stack:**

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Elasticsearch](https://img.shields.io/badge/Elasticsearch-005571?style=for-the-badge&logo=elasticsearch&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/RabbitMQ-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)

## Architecture & Core Logic
* **Coordinator Service (Go):** Handles API requests, manages tasks via RabbitMQ, and tracks real-time status in Redis.
* **Scrapers (Go):** High-performance workers using the `Colly` library to extract job listings from various platforms.
* **Data Processor (Go):** Processes incoming data streams, normalizes information, and persists it to storage.
* **Search Engine:** Powered by Elasticsearch using `MatchQuery` with `Fuzziness: "AUTO"` to handle typos in location and job titles.

### High-level flow:
`Client` → `CoordinatorService` → `RabbitMQ` → `Scrapers` → `DataService` → `PostgreSQL + Elasticsearch` → `CoordinatorService` → `Search Results`

<p align="center">
  <img height="600" alt="picture" src="https://github.com/user-attachments/assets/7ac8ba10-0d69-477e-b20e-828baf422eff" />
  <img height="600" alt="picture" src="https://github.com/user-attachments/assets/608c22d5-0845-4359-82a2-74127e759d13" />
</p>

---

## API Reference

### Endpoints

| Type | Endpoint         | Description                                                                                                |
|------|------------------|------------------------------------------------------------------------------------------------------------|
| `POST` | /tasks           | Initiates a new background scraping process.                                                               |
| `GET`  | /tasks/{task_id} | Provide a non-blocking, real-time status update for the background work.                                   |
| `GET`  | /sources         | Provides a list of supported platforms.                                                                    |
| `GET`  | /vacancies       | Provides a clean, stable, and fast interface for the end-user to interact with the collected intelligence. |

### Query Parameters (`GET /vacancies`)

| Parameter    |   Type  | Required | Description                  |
|--------------|:-------:|-------------|------------------------------|
| query        |  string |     not     | Search query (e.g. "Golang") |
| requirements |  string |     not     | Job requirements             |
| location     |  string |     not     | City or region to search     |
| limit        | integer |     not     | Number of results per page   |
| offset       | integer |     not     | Offset for pagination        |

---

## Data Format (JSON)

### JSON (Request/Response Body)
Request Body (POST /tasks):
```json
{
  "keywords": ["junior", "python"],
  "sources": ["dou.ua"]
}
```
Response Body(POST /tasks):
```json
{
    "task_id": "74385192-4be2-45df-8a5b-f573bb341120",
    "status": "created",
    "created_at": "2026-02-02T16:57:32Z"
}
```

Request Body (GET /vacancies):
| Parameter    |   Value  |
|--------------|:-------:|
| query        |  golang |
| requirements |  "" | 
| location     |  "Київ" | 
| limit        | 10 | 
| offset       | 0 | 

Response Body(POST /tasks):
```json
{
    "items": [
        {
            "id": "e41a7823-7ace-4752-8650-5cfe379b47e1",
            "title": "Junior Golang Developer",
            "company": "Cossack Labs",
            "location": "Київ, Львів, віддалено",
            "description": "This position is open only to..",
            "link": "https://jobs.dou.ua/companies/cossack-labs/vacancies/...",
            "about": "",
            "requirements": "• At least 1.5 years..."
        }
    ],
    "total": 1
}
```

Response Body(GET /sources):
```json
{
    "sources": [
        {
            "id": "99c32bcf-2c64-4ddc-a72d-13ed6cd45dd9",
            "name": "dou.ua"
        }
    ],
    "total": 1
}
```

---

## Environment Variables (.env)
|     **Parameter**     |                          **Example Value**                          |
|:---------------------:|:-------------------------------------------------------------------:|
| POSTGRES_USER         | postgres                                                            |
| POSTGRES_PASSWORD     | secure_db_pass_2026                                                 |
| POSTGRES_DB           | service_db                                                          |
| DATA_PROCESSOR_DB_URL | postgres://example:example@postgres:5432/example_db?sslmode=disable |
| RABBITMQ_USER         | guest                                                               |
| RABBITMQ_PASS         | guest                                                               |
| RABBIT_URL            | amqp://guest:guest@rabbitmq:5672/                                   |
| ELASTICSEARCH_HOST    | elasticsearch                                                       |
| ELASTICSEARCH_PORT    | 9200                                                                |
| ELASTICSEARCH_URL     | http://elasticsearch:9200                                           |
| REDIS_HOST            | redis                                                               |
| REDIS_PORT            | 6379                                                                |
| REDIS_PASSWORD        | secure_db_pass_2026                                                 |
| HTTP_PORT             | 8080                                                                |

---

## Quick Start

### Prerequisites

- Docker + Docker Compose
- Git
- Go 

### Run with Docker Compose (recommended)

```bash
# 1. Clone
git clone https://github.com/ReilEgor/Vaca
cd Vaca/deployments


# 2. Fill .env with Environment Variables
touch .env

# 3. Run with Docker Compose
docker-compose up --build
```
Wait ~30–90 seconds until services are ready.

---

## License

Distributed under the MIT License. See `LICENSE` for more information.

**Developed by [YehorReil](https://github.com/ReilEgor)**
