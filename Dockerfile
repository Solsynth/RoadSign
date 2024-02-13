# Building Backend
FROM rust:alpine as roadsign-server

RUN apk add nodejs npm

WORKDIR /source
COPY . .
RUN apk add libressl-dev
RUN cargo build --release

# Runtime
FROM alpine:latest

COPY --from=roadsign-server /source/target/release/roadsign /roadsign/server

EXPOSE 81

CMD ["/roadsign/server"]