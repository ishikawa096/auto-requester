FROM golang:1.23-alpine3.19 AS development

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["go", "run", "main.go"]

FROM golang:1.23-alpine3.19 AS build

WORKDIR /app

RUN apk add --no-cache upx

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -buildvcs=false -trimpath -ldflags '-w -s' -o main .
RUN upx main

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=build /app/main /bin/main

ENTRYPOINT ["/bin/main"]
