FROM rust:1.80.1 AS builder

WORKDIR /app

COPY . .

RUN cargo build --release

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update -y \
    && apt-get install -y --no-install-recommends openssl ca-certificates \
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/release/ml-service /app/ml-service
COPY ./models/model_0.0.1.onnx /app/models/model_0.0.1.onnx
COPY ./preprocess/scaler_tensors.csv /app/preprocess/scaler_tensors.csv
COPY configuration configuration

EXPOSE 7216

CMD ["./ml-service"]
