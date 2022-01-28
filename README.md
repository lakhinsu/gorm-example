# Example for using ORM in Golang with GORM
This repository contains the example code for using ORM in golang. I have also written a [blog]() for the same, feel free to check it out. 
This repository also contains Dockerfile and docker-compose.yaml files for deploying these example using Docker.

## Features
This repo contains example code for following points,

- Creating a GORM model.
- Creating a connection pool with the database.
- Auto migrate the models to create the desired tables automatically in the database.
- Example code for CRUD operations.

## Run locally
### API Service
Follow these steps to run the API service locally,
At the repository root, execute the following commands.
- Get dependencies.
`go get .`
- Set the environment variables - use the sample .env file provided in the repository.
- Execute the `go run main.go` command to start the service.
> Note: This assumes you already have a database instanse running. If not, use the following guide to run both API service and a database instance in Docker.
## Running in Docker
- Running in Docker is very simple, the repo includes a `docker-compose.yaml` file.
- At the repository root, execute `docker-compose up` command to deploy the Postgres database instance and the API service.

> The repository includes a Insomnia collection, import it into the Insomnia REST client to play with the REST service.
