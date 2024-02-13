# Building Backend
FROM rust:alpine as roadsign-server

RUN apk add libressl-dev build-base

WORKDIR /source
COPY . .
ENV RUSTFLAGS="-C target-feature=-crt-static"
RUN cargo build --release

# Runtime
FROM alpine:latest

COPY --from=roadsign-server /source/target/release/roadsign /roadsign/server

EXPOSE 81

CMD ["/roadsign/server"]