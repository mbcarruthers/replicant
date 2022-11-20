## Replicant

A boilerplate for containerized go/gin projects.


Todo: 
- optional: If possible, find a way to run the makefile within the docker container for building. Not necessary but when
  make is installed in replcant.Dockerfile (`RUN apk add make`) my makefile does not recognize the Alpine operating
  system. Therefore, in the Makefile's OS rules Alpine will have to be added as an operating system.

- Attach an Envoy container(s) for everything from load-balancing, HTTPs proxy, HTTP2 proxy, TLS , whatever.
  Just get it involved in some way and go from there. i.e - make a dir in the container dir container the envoy Dockerfile
 as well an the envoy.yaml file and add it to the docker-compose.yaml
- Attach cockroachdb to the server and have it 'parametized' to connect to that service in the docker-compose.yaml
