# Building Backend
FROM golang:alpine as roadsign-server

RUN apk add nodejs npm

WORKDIR /source
COPY . .
WORKDIR /source/pkg/sideload/view
RUN npm install
RUN npm run build
WORKDIR /source
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs -o /dist ./pkg/cmd/server.rs/main.go

# Runtime
FROM golang:alpine

COPY --from=roadsign-server /dist /roadsign/server

EXPOSE 81

CMD ["/roadsign/server"]