FROM golang AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/app/main.go
FROM ubuntu as runner

WORKDIR /app

COPY --from=build /app/main .
CMD [ "./main" ]