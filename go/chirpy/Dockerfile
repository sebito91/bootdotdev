# start of with debain OS
FROM debian:stable-slim

# COPY source destination
COPY chirpy /bin/chirpy
COPY .env /.env
COPY index.html /index.html

# bind to port 8080
ENV PORT 8080

# run the server on startup
CMD ["/bin/chirpy"]

