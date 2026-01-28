# Coordinator Service

The **Coordinator Service** is the central orchestrator and the primary **entry point** for the entire microservices ecosystem. It acts as an **API Gateway** and a task manager, ensuring seamless communication between the client-side and internal worker services.

### API Reference
| Type | Endpoint         | Description                                                                                                |
|------|------------------|------------------------------------------------------------------------------------------------------------|
| `POST` | /tasks           | Initiates a new background scraping process.                                                               |
| `GET`  | /tasks/{task_id} | Provide a non-blocking, real-time status update for the background work.                                   |
| `GET`  | /sources         | Provides a list of supported platforms.                                                                    |
| `GET`  | /vacancies       | Provides a clean, stable, and fast interface for the end-user to interact with the collected intelligence. |
<p align="center">
<img height="480" alt="picture" src="https://github.com/user-attachments/assets/608c22d5-0845-4359-82a2-74127e759d13" />
</p>
