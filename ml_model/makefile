.PHONY: download-dataset
download-dataset:
	@echo "Download dataset"
	@rye run kaggle datasets download -p ../data/ -d stephanmatzka/predictive-maintenance-dataset-ai4i-2020
	@echo "Unzip dataset"
	@unzip -o ../data/predictive-maintenance-dataset-ai4i-2020 -d ../data/
