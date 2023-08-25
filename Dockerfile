FROM golang:1.20 as build-stage
RUN mkdir /app
ADD . /app/
WORKDIR /app
ENV GOBIN /app
RUN CGO_ENABLED=0 GOOS=linux go install github.com/hoveychen/slime@latest

FROM gcr.io/distroless/base-debian11 AS build-release-stage
COPY --from=build-stage /app/slime /app/slime
WORKDIR /app
ENTRYPOINT ["/app/slime"]
