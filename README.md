# API-CI-CD
## Author
### Rain Praks
## Project Overview
This project is a Golang-based API that acts as a proxy to the Pipedrive API, forwarding requests while logging all interactions. It includes monitoring via a `/metrics` endpoint and follows a CI/CD pipeline using GitHub Actions.
## Features
- **API Endpoints**
  - `GET /deals` - Retrieves all deals from the Pipedrive API.
  - `POST /deals` - Adds a new deal.
  - `PUT /deals/{id}` - Updates an existing deal.
  - `DELETE /deals/{id}` - Deletes a deal.
  - `GET /metrics` - Provides request duration and latency metrics.
- **Instrumentation**
  - Logs all requests and responses in a structured format.
  - Captures request durations and latencies for monitoring.
- **CI/CD Implementation**
  - **CI:** Runs tests, linting, and static analysis on every pull request.
  - **CD:** Deploys automatically when changes are merged into `main`.
- **Docker Support**
  - The API is containerized with a multi-stage build process for efficiency.
## Components
### API
- **Handlers**: Implements business logic for forwarding API requests.
- **Middleware**: Logs requests and gathers metrics.
- **Main**: Initializes the router and starts the server.
### CI/CD
- **CI Workflow** (`.github/workflows/ci.yml`)
  - Runs tests, formatting checks, and linting on PRs to `develop` and `main`.
- **CD Workflow** (`.github/workflows/cd.yml`)
  - Logs "Deployed!" into GitHub Action UI when a PR is merged into `main`.
### Docker
- **Dockerfile**
  - Uses a multi-stage build to create a minimal, production-ready image.
  - Exposes port `8080` for API requests.
## How to Run
### Running locally
- create .env file to the root directory:
```
PIPEDRIVE_API_TOKEN=your_api_token
PIPEDRIVE_API_URL=https://api.pipedrive.com/v1/deals
```
- run commands:
    - ```go mod tidy```
    -  ```go run . ```
### Running with Docker
- create .env file like in previous step.
- run commands in root directory:
    - To build image:```docker build -t [image name] .```
    - Check images: ```docker images```
    - Start the container: ```docker run -d -p 8080:8080 --name [image name] [container name]```
    - check if container runs: ```docker ps```
    - to stop the container ```docker stop [container name]```
    - to delete the container: ```docker rm [container name or id]```
    - to delete the image: ```docker rmi [image name or id]```
## How to Test Endpoints
### Using Postman
1. Open Postman and create a new request.
2. Set the request method and endpoint URL:
   - **GET /deals**: `http://localhost:8080/deals`
   - **POST /deals**: `http://localhost:8080/deals`
   - **PUT /deals/{id}**: `http://localhost:8080/deals/{id}`
   - **DELETE /deals/{id}**: `http://localhost:8080/deals/{id}`
   - **GET /metrics**: `http://localhost:8080/metrics`
3. Set request headers:
   - `Content-Type: application/json`
4. For `POST` and `PUT`, include a JSON body, for example:
```json
{
  "title": "New Deal",
  "value": 1000
}
```
5. Send the request and verify the response.
### Using curl
```sh
# Get all deals
curl http://localhost:8080/deals

# Create a new deal
curl -X POST http://localhost:8080/deals -H "Content-Type: application/json" -d '{"title": "New Deal", "value": 1000}'

# Update an existing deal
curl -X PUT http://localhost:8080/deals/[id] -H "Content-Type: application/json" -d '{"title": "Updated Deal", "value": 1500}'

# Delete a deal
curl -X DELETE http://localhost:8080/deals/[id]

# Get metrics
curl http://localhost:8080/metrics
```
## Notes
- The API uses environment variables for Pipedrive credentials.
- Tests use a mock server to simulate Pipedrive API responses.
- /metrics does not provide latency measurements because there is no access to a load balancer, making it impossible to accurately determine request latency.
- GitHub Actions handles automated testing and deployment.
