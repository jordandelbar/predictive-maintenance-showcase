version ?= 0.0.1-dev

.PHONY: dataset/download
dataset/download:
	@echo "Download dataset"
	@rye run kaggle datasets download -p ../data/ -d nphantawee/pump-sensor-data

.PHONY: dataset/unzip
dataset/unzip:
	@echo "Unzip dataset"
	@unzip -o ../data/pump-sensor-data.zip -d ../data/

.PHONY: dataset/all
dataset/all: dataset/download dataset/unzip

.PHONY: model/test
model/test:
	@echo "Test"
	@rye run pytest --cov=src tests/

.PHONY: model/train
model/train:
	@echo "Train model"
	@rye run python src/ml_model/main.py

.PHONY: model/build
model/build:
	@echo "Build model version ${version}"
	@rye run bentoml build --containerize --version ${version} -f bentofile.yaml ./src/ml_model/
