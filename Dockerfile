# Building Backend
FROM rust:alpine as roadsign-server

RUN apk add libressl-dev build-base

WORKDIR /source
COPY . .
ENV RUSTFLAGS="-C target-feature=-crt-static"
RUN cargo build --release

EXPOSE 81

CMD ["/source/target/release/roadsign"]