from typing import Tuple, List

import bentoml
import numpy as np
import matplotlib.pyplot as plt
import torch
import torch.nn as nn
import torch.optim as optim

from loguru import logger
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import MinMaxScaler

import ml_model


def split_data(
    x: np.ndarray,
    y: np.ndarray,
    test_size: float = 0.3,
    random_state: int = 42,
) -> Tuple[np.ndarray, np.ndarray, np.ndarray, np.ndarray]:
    x_train, x_test, y_train, y_test = train_test_split(
        x, y, test_size=test_size, stratify=y, random_state=random_state
    )
    logger.info(f"x train shape: {x_train.shape}")
    logger.info(f"x test shape: {x_test.shape}")
    logger.info(f"y train shape: {y_train.shape}")
    logger.info(f"y test shape: {y_test.shape}")
    return x_train, x_test, y_train, y_test


def train_model(
    x_train: np.ndarray, y_train: np.ndarray, epochs: int = 50
) -> ml_model.AutoEncoder:
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

    scaler = MinMaxScaler()
    scaler.fit(x_train)

    # We use the Scikit-learn MinMaxScaler and instantiate our autoencoder with
    # the min and max features vectors
    autoencoder = ml_model.AutoEncoder(
        input_dim=52,
        hidden_dim=24,
        min_tensor=scaler.data_min_,
        max_tensor=scaler.data_max_,
        min_scaling_range=0.0,
        max_scaling_range=1.0,
    ).to(device)
    criterion = nn.MSELoss()
    optimizer = optim.Adam(autoencoder.parameters(), lr=0.001)
    x_train = autoencoder.min_max_scaling(torch.tensor(x_train, dtype=torch.float32))

    model, training_loss_list = ml_model.train_model(
        x_train=x_train,
        optimizer=optimizer,
        autoencoder=autoencoder,
        criterion=criterion,
        device=device,
        epochs=epochs,
    )
    plt.plot(range(len(training_loss_list)), training_loss_list)
    plt.savefig("./output/training_loss.png")
    return model


def evaluate_model(
    x_test: np.ndarray,
    model: ml_model.AutoEncoder,
) -> List[float]:
    x_test = torch.tensor(x_test, dtype=torch.float32)
    reconstruction_errors = model.predict(x_test)
    logger.info(f"reconstruction errors: {reconstruction_errors}")

    # TODO: verify if it is batchable and the threshold
    signatures = {"predict": {"batchable": True}}
    metadata = {"best_threshold": 0.1}
    saved_model = bentoml.pytorch.save_model(
        "pm_autoencoder",
        model,
        signatures=signatures,
        metadata=metadata,
        external_modules=[ml_model],
    )
    logger.info(f"Model saved: {saved_model}")

    return reconstruction_errors.numpy()
