FROM golang:1.15 AS build
WORKDIR /workspaces
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN cd tester; go build -o tester main.go
RUN cd observer; go build -o observer main.go

FROM debian:buster AS tester
WORKDIR /app
COPY --from=build /workspaces/tester/tester /app/
ENTRYPOINT ./tester

FROM debian:buster AS observer
WORKDIR /app
COPY --from=build /workspaces/observer/observer /app/
ENTRYPOINT ./observer
