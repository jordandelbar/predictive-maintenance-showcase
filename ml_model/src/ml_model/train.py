from pathlib import Path
from typing import Tuple

import polars as pl
import numpy as np
import torch
import torch.nn as nn
import torch.optim as optim

from loguru import logger
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import MinMaxScaler
from sklearn.metrics import f1_score, recall_score, precision_score, confusion_matrix

import ml_model


def load_csv(filename: str) -> pl.DataFrame:
    logger.info(f"Loading data from {filename}")
    data = pl.read_csv(f"{Path(__file__).parents[3]}/data/{filename}.csv")
    return data


def preprocess_data(data: pl.DataFrame) -> Tuple[np.ndarray, np.ndarray]:
    training_data = data.select(
        [
            "Air temperature [K]",
            "Process temperature [K]",
            "Rotational speed [rpm]",
            "Torque [Nm]",
            "Tool wear [min]",
            "Machine failure",
        ]
    )
    x_df = training_data.select(pl.exclude("Machine failure"))
    y_df = training_data.select(pl.col("Machine failure"))

    logger.info(f"x dataframe shape: {x_df.shape}")
    logger.info(f"y dataframe shape: {y_df.shape}")
    x_array = x_df.to_numpy()
    y_array = y_df.to_numpy().squeeze()
    return x_array, y_array


def split_data(x: np.ndarray, y: np.ndarray, test_size: float = 0.3):
    x_train, x_test, y_train, y_test = train_test_split(
        x, y, test_size=test_size, stratify=y
    )
    logger.info(f"x train shape: {x_train.shape}")
    logger.info(f"x test shape: {x_test.shape}")
    logger.info(f"y train shape: {y_train.shape}")
    logger.info(f"y test shape: {y_test.shape}")
    return x_train, x_test, y_train, y_test


def train_model(
    x_train: np.ndarray, y_train: np.ndarray
) -> Tuple[ml_model.AutoEncoder, float]:
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

    x_normal_events, _ = ml_model.return_normal_events_sets(
        x_train=x_train, y_train=y_train
    )

    scaler = MinMaxScaler()
    scaler.fit(x_train)

    # We use the Scikit-learn MinMaxScaler and instantiate our autoencoder with
    # the min and max features vectors
    autoencoder = ml_model.AutoEncoder(
        input_dim=5,
        hidden_dim=20,
        min_tensor=scaler.data_min_,
        max_tensor=scaler.data_max_,
        min_scaling_range=0.0,
        max_scaling_range=1.0,
    ).to(device)
    criterion = nn.MSELoss()
    optimizer = optim.Adam(autoencoder.parameters(), lr=0.001)
    x_ne_scaled = autoencoder.min_max_scaling(
        torch.tensor(x_normal_events, dtype=torch.float32)
    )

    model = ml_model.train(
        x_ne_scaled=x_ne_scaled,
        optimizer=optimizer,
        autoencoder=autoencoder,
        criterion=criterion,
        device=device,
        epochs=50,
    )
    best_threshold = ml_model.compute_best_threshold(x_train, y_train, model, 1.5)
    logger.info(f"{best_threshold=}")
    return model, best_threshold


def evaluate_model(
    x_test: np.ndarray,
    y_test: np.ndarray,
    model: ml_model.AutoEncoder,
    threshold: float,
) -> None:
    x_test = torch.tensor(x_test, dtype=torch.float32)
    reconstruction_errors = model.predict(x_test)
    y_pred = np.where(reconstruction_errors >= threshold, 1, 0)
    f1 = f1_score(y_true=y_test, y_pred=y_pred)
    recall = recall_score(y_true=y_test, y_pred=y_pred)
    precision = precision_score(y_true=y_test, y_pred=y_pred)

    logger.info(f"{f1=}")
    logger.info(f"{recall=}")
    logger.info(f"{precision=}")
    logger.info(f"{confusion_matrix(y_true=y_test, y_pred=y_pred)}")

    _ = {"predict": {"batchable": True}}
    _ = {"best_threshold": threshold}
    # saved_model = bentoml.pytorch.save_model(
    #     "pm_autoencoder",
    #     model,
    #     signatures=signatures,
    #     metadata=metadata,
    #     external_modules=[ml_model],
    # )
    # logger.info(f"Model saved: {saved_model}")
