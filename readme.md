# DIS PROJECT 2024YOLO

- Golang + htmx

## Running the Project

Do this to get a copy of the project up and running on your local machine for demonstration purposes.

### Prerequisites

1. **Locally** If running locally, get [Golang](https://go.dev/dl/) and [PostgreSQL](https://www.postgresql.org/)
2. **If using Docker (preferred, if your not developing)**: Download and install Docker from [Docker Desktop](https://www.docker.com/products/docker-desktop)

### Installation and running the project

#### Option 1: Running locally

1. Clone the repository:

   ```sh
   git clone https://github.com/kongpng/DIS_PROJECT.git
   cd DIS_PROJECT
   ```

2. Setup postgres, and edit the .env file to your postgres database settings.

3. Run the database scripts (TBA)

4. Run

```sh
   go run project
```

5. Go to localhost:8080 on your browser.

#### Option: Using Docker

1. Clone the repository:

   ```sh
   git clone https://github.com/kongpng/DIS_PROJECT.git
   cd DIS_PROJECT
   ```

2. Run docker compose:

   ```
   docker compose up --build -d, if you already built once, then just docker compose up
   ```

3. go to localhost:8080
