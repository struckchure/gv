FROM golang:1.24.1-alpine AS builder

WORKDIR /code

COPY . .

ENV GV_MODE=build

RUN go mod tidy
RUN go run .
RUN go build -o ./main .

FROM istio/distroless:1.26-2025-04-24T19-01-35

WORKDIR /code

COPY --from=builder /code/dist dist
COPY --from=builder /code/main .

EXPOSE 3000

CMD [ "./main" ]
