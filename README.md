## Replicant

A boilerplate for containerized go/gin projects with a single-node cockroachDB connection

The Go container creates a 'test database' at startup to show it connects successfully by creating a trivial database+table, inserting an element, and querying it. 

The Go container is started by a multistage Dockerfile and will not run until the healthcheck for the database container shows that it is ready for connection. The multistage
dockerfile builds the image inside of an alpine container with access to the libaries for compilation and moves it to an base alpine image for hosting.

The gin server comes along with a graceful shutdown. 

The cockroach docker container contains no volume for persistence and thus is completely removed when 'docker-compose down' is ran.
