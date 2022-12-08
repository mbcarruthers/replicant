## Replicant

A boilerplate for containerized go/gin projects with a single-node cockroachDB connection

The Go container creates a 'test database' at startup to show it connects successfully by creating a trivial database+table, inserting an element, and querying it. 

The Go container is started by a multistage Dockerfile and will not run until the healthcheck for the database container shows that it is ready for connection. The multistage
dockerfile builds the image inside of an alpine container with access to the libaries for compilation and moves it to an base alpine image for hosting.

The gin server comes along with a graceful shutdown. 

The cockroach docker container contains no volume for persistence and thus is completely removed when 'docker-compose down' is ran.

App exposes 2 simple endpoints available outside the container.

1. localhost:8000/api/echo 
    - simple echo response to to really whatever is sent in the body to the endpoint. If you send 'echo' it response with 'echo'. If you send a sonnet by shakespear, it replys it back.

2. localhost:8000/api/data
    - response with the single content-string of the temporary database, which is 'Decker'.

And Will end with a graceful shutdown.
