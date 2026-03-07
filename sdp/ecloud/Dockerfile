  FROM --platform=linux/amd64 rust:1.91 AS builder
  WORKDIR /app
  COPY . .
  RUN cargo build --release -p sdp-relayer

  FROM --platform=linux/amd64 debian:bookworm-slim
  RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
  WORKDIR /app
  COPY --from=builder /app/target/release/sdp-relayer .
  EXPOSE 3000
  CMD ["./sdp-relayer"]