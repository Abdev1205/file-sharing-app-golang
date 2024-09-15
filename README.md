## Installation Guide

### Prerequisites

Before you begin, ensure you have the following installed:

- Docker
- Docker Compose

### Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Abdev1205/21BCE11045_Backend.git
   ```

2. **Build and start the application:**

   Use Docker Compose to build and start the application along with its dependencies.

   ```bash
   docker-compose up --build
   ```

   This command will:

   - Build the Docker images for the Go application, PostgreSQL, and Redis and then project will start

3. **Verify the setup:**

   Once the containers are up, you can verify that the application is running by accessing the server's endpoint. By default, it should be available at `http://localhost:8080`.

## Project Structure

- `/server`: Contains the entry point for the server.
  - `main.go`: Starts and configures the server.
- `/internal`: non-exportable packages.
  - `/adapters`: Infrastructure layer (database, file handling, caching).
    - `/database`: PostgreSQL database.
    - `/cache`: Redis caching.
  - `/application`: Application layer.
    - `/file_service`: File service logic (Upload, Retrieve, etc.).
    - `/auth_service`: JWT authentication logic (Login, Register, etc.).
- `/pkg`: Shared utilities.
  - `/middleware`: Middleware for JWT authentication and other purposes.
  - `/config`: Configuration management (DB, Redis, JWT, etc.).
- `docker-compose.yml` and `Dockerfile`: Configuration for Docker Compose and Docker build.
- `go.mod`: Go module dependencies.




## Demo
[Postman Api Endpoints](https://lively-comet-969560.postman.co/workspace/My-Workspace~98fbf43a-20f9-4cae-8b31-6b9d4ed9a210/collection/23044745-45ab4c8e-df12-4561-a71d-1ea150dc13f6?action=share&creator=23044745) 

## Project Starting
![image](https://github.com/user-attachments/assets/1abc49a2-5096-42d8-8954-066b6d5152be)


## Register User
![image](https://github.com/user-attachments/assets/aa00e3c8-8de5-4312-a638-30054368ed45)

## Login User

![image](https://github.com/user-attachments/assets/0132b90c-c94f-40dd-8a3e-77da77cb19e1)

- ### also getting jwt cookies 

![image](https://github.com/user-attachments/assets/694e5485-9985-4db5-8f26-63435ed1c5b3)

## File Upload 
- ### Authenticated Route requires Authorization Bearer Token

![image](https://github.com/user-attachments/assets/2f921780-4ed1-48a7-8693-d19100c17b4f)

- ### After passing Token files upload successfully locally in upload folder

![image](https://github.com/user-attachments/assets/1f9841f6-8d76-4ed3-83ec-689fea32fe4f)

- ### We can see our files from /upload routes 

![image](https://github.com/user-attachments/assets/05e575ba-6a42-42cd-b73a-94ab4c9ec7ac)

- ### Getting All files 

![image](https://github.com/user-attachments/assets/05635a98-1b16-45d8-b68a-0849d767f22d)

- ### Getting shareable link by passing file id 

![image](https://github.com/user-attachments/assets/04f25e91-b067-41b3-96b7-b8c39c074950)

- ### Searching Files based on Query Params

![image](https://github.com/user-attachments/assets/0a542818-c31f-4009-a806-831987786ee4)



