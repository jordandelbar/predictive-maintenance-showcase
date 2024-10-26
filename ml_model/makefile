version ?= 0.0.1
epochs ?= 50

.PHONY: model/test
model/test:
	@uv run pytest --cov=src tests/

.PHONY: model/train
model/train:
	@uv run python src/ml_model/main.py --epochs=$(epochs) --version=$(version)
