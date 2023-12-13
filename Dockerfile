# Building Backend
FROM golang:alpine as roadsign-server

WORKDIR /source
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs -o /dist ./pkg/cmd/server/main.go

# Runtime
FROM golang:alpine

COPY --from=roadsign-server /dist /roadsign/server

EXPOSE 81

CMD ["/roadsign/server"]