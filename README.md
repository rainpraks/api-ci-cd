Software Engineering Intern DevOps
home challenge
Your task is split into five parts, complete as many as you can. Email us the public
GitHub repository link once you have finished. It is preferred to commit as you develop,
instead of having one single commit.
Part I - API
Write a simple application in any language (Typescript or Golang preferred), with three
endpoints, that all forward requests to the public Pipedrive API
(https://developers.pipedrive.com/docs/api/v1/Deals). Use api-tokens.
●
●
●
GET /deals - gets all deals
POST /deals - adds a new deal
PUT /deals - updates an already existing deal
Test as much as it is reasonable.
Part II - Instrumentation
Log everything that happens, to console, in a clean and succinct way.
Add a GET /metrics endpoint that will return request metrics (such as mean request
duration and latency) about all endpoints.
Part III - CI
Use GitHub Actions to create a simple test-on-commit workflow that will run all your
tests, and linting, for every commit pushed to a pull-request.
Part IV - CD
Use GitHub Actions once again to create a workflow that runs only when a pull-request
is merged to master. This workflow should do nothing but log “Deployed!”
Part V - Reproducibility
Create a README with clear instructions on how to run the application locally,
preferably (strongly) with Docker. Describe every part that you have done.
→ Submitting the challenge
Once you are ready to submit, we would like you to share your public GitHub repository
link with Ragnar (ragnar.paide@pipedrive.com), Martin (martin.aun@pipedrive.com),
Karoline (karoline.ansberg@pipedrive.com) and Keiri (keiri.tohvrel@pipedrive.com).
Good luck!