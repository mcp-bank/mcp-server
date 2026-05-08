FROM golang:1.26-alpine AS build
WORKDIR /www/app
COPY go.mod go.sum ./
RUN ["go", "mod", "download"]
COPY . .
RUN ["go", "build", "-o", "/www/app/main", "./cmd/"]
FROM alpine
WORKDIR /www/app
COPY --from=build /www/app/main .
CMD ["./main"]