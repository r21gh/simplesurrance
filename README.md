# Simplesurrance technical interview
Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total 
number of requests that it has received during the previous 60 seconds (moving window). 
The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

## Development
In this project I used the standard library to create a simple HTTP server.

## How to run
In order to run the server, you need to install the Go language.

```go run main.go```

## What is the output
The server should return the correct number of requests during the previous 60 seconds.

## What is the input
The server should accept any request.

## How to test
The test file for the server is under the services folder.
```go run main.go```

## How to run using Docker
At first go to the main directory of the project.

```docker build -t simplesurrance .```

Then run the container

```docker run --rm -d -p 9000:9000 simplesurrance```
