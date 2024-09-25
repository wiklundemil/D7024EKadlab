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


FROM golang:latest

WORKDIR /app

# Copy the local code from the Go directory to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose a port (if needed)
EXPOSE 3000

# Define the command to run the executable
CMD ["./main"]