.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: download/dataset
download/dataset:
	@echo "Download dataset"
	@rye run kaggle datasets download -p ../data/ -d stephanmatzka/predictive-maintenance-dataset-ai4i-2020
	@echo "Unzip dataset"
	@unzip -o ../data/predictive-maintenance-dataset-ai4i-2020 -d ../data/

.PHONY: model/train
model/train:
	@echo "Train model"
	@rye run python src/ml_model/main.py

.PHONY: model/test
model/test:
	@echo "Test"
	@rye run pytest --cov=src tests/
