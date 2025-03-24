# Apollo

## Overview
Apollo is a backend service for managing various backend operations. This service is built using Go and relies on PostgreSQL for data storage. The project is structured to be easily set up and run using a Makefile, which simplifies the build and deployment process.

## Requirements
Before you begin, ensure you have the following installed on your system:
1. **Go 1.23.0 or later**: The service is written in Go, so you need to have Go installed. You can download and install it from the [official Go website](https://go.dev/dl/).
2. **PostgreSQL**: The service uses PostgreSQL as its database. Make sure you have PostgreSQL installed and running. You can download it from the [official PostgreSQL website](https://www.postgresql.org/download/).
3. **CMake**: CMake is required for building certain dependencies. You can install it from the [official CMake website](https://cmake.org/download/).

## Installation
### 1. **Clone the Repository:** 
Open your terminal and run the following command to clone the repository:</br>
```bash 
  git clone https://github.com/winartodev/apollo.git
```

### 2. Navigate to the Project Directory
Change to the directory where the repository was cloned:
```bash 
  cd apollo
```

3. Set Up the Database
Ensure PostgreSQL is running and create a new database for Apollo. Update the database connection details in the [configuration file](https://github.com/winartodev/apollo/blob/main/core/files/apollo.dev.yaml.template)


## Configuration
Apollo uses a [configuration file](https://github.com/winartodev/apollo/blob/main/core/files/apollo.dev.yaml.template) to manage settings such as database connections, server ports and so on. Below is an example configuration:
### 1. Copy Template Files
copy template file `apollo.dev.yaml.template` to `apollo.dev.yaml`
```yaml
app:
  name: Apollo
  port:
    http: 8989
database:
  driver: postgres
  host: 127.0.0.1
  port: 5432
  name: apollo_db
  username: apollo_user
  password: apollo
  sslMode: disable
  defaultMaxConn: 200
  defaultIdleConn: 20
  connMaxLifetime: 10 # in minutes
  connMaxIdleTime: 5 # in minutes
auth:
  apiKey: 123
  jwt:
    accessToken:
      secretKey: # <your access token secret key>
    refreshToken:
      secretKey: # <your refresh token secret key>
```

## Usage
Once the service is running, you can interact with it via HTTP requests. Below are some example endpoints:

### Health Check
```bash
curl http://localhost:8989/api/healthz
```

