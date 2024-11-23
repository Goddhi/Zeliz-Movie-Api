# Zeliz Movie API - Golang Backend Project
### Overview
The Movie API is a RESTful service developed using Golang that supports advanced CRUD operations, designed for high scalability, security, and performance. This API allows clients to interact with movie-related data efficiently while ensuring data integrity and system stability.

### Features
Optimistic Concurrency Control: Prevents conflicts during simultaneous updates by implementing version-based checks to maintain data consistency.
Partial Updates (PATCH): Enables clients to modify only specific fields of a resource, optimizing data handling and reducing unnecessary payloads.
Rate Limiting: Both IP-based and global rate limiting are implemented to manage traffic, ensuring fair usage and preventing overloads.
Data Storage: Utilizes PostgreSQL for efficient database management, schema design, and optimized query performance.
Microservices Architecture: The API follows a modular approach, ensuring scalability and maintainability through independent services, adhering to SOLID design principles and clean code practices.
Technologies Used
**Golang:** The backend is developed using Golang to leverage its concurrency handling and performance efficiency for high-load applications.
**PostgreSQL:** Used for data storage, ensuring efficient database management with well-structured schemas and optimized queries.
**Microservices:** A microservices architecture was employed to ensure each service is independent, scalable, and easy to maintain.
**GitHub Actions:** For automating the CI/CD pipeline, ensuring continuous integration and seamless deployment.
**Docker:** The application is containerized using Docker, ensuring consistent environments across development, testing, and production stages.

### Key Features and Implementations
##### Optimistic Concurrency Control:
Ensures that updates to the same data by multiple clients do not cause conflicts. A versioning system was implemented to compare and check if the data has been modified before applying updates.

#### Updates (PATCH):
Allows clients to send only the modified fields in a request, reducing the payload size and increasing the efficiency of data updates.

#### Rate Limiting:
Implements both IP-based and global rate limiting to control the number of requests and ensure fair usage, enhancing system stability and preventing abuse.

#### Database with PostgreSQL:
Designed a robust schema in PostgreSQL to efficiently handle movie data, including genres, cast, and ratings, and optimized queries for quick data retrieval.

#### Microservices Architecture:
The API is built on a microservices architecture, where each service is isolated and can be developed, deployed, and scaled independently. This ensures better maintainability and easier updates.
