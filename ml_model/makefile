version ?= 0.0.1
epochs ?= 50

.PHONY: model/test
model/test:
	@rye run pytest --cov=src tests/

.PHONY: model/train
model/train:
	@rye run python src/ml_model/main.py --epochs=$(epochs) --version=$(version)
