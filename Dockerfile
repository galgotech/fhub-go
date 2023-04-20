FROM golang:1.20.3 as build
WORKDIR /app
COPY . /app/
RUN go build -o /usr/local/bin/fhub-gencode ./cmd/gencode
RUN go build -o /usr/local/bin/fhub-rest ./cmd/rest

FROM golang:1.20.3 as runtime
WORKDIR /app
COPY --from=build /usr/local/bin/fhub-gencode /usr/local/bin/
COPY --from=build /usr/local/bin/fhub-rest /usr/local/bin/
