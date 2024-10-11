FROM alpine:latest

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab


# Stage 1: Build the Go binary
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o kademlia_app main.go

# Stage 2: Create a small runtime image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/kademlia_app .
EXPOSE 3000
CMD ["./kademlia_app"]
