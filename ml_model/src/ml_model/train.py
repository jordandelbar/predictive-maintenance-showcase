from pathlib import Path

import matplotlib.pyplot as plt
import numpy as np
import polars as pl
import torch
import torch.nn as nn
import torch.optim as optim

from loguru import logger
from sklearn.preprocessing import MinMaxScaler
from torch.utils.data import DataLoader, TensorDataset

from ml_model.model import AutoEncoder
from ml_model.utils import remove_graph


def train_model(
    x_train: np.ndarray,
    epochs: int = 25,
    patience: int = 5,
    version: str = "0.0.1",
) -> AutoEncoder:
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

    scaler = MinMaxScaler()
    scaler.fit(x_train)

    min_array = scaler.data_min_
    max_array = scaler.data_max_

    pl.DataFrame({"min_array": min_array, "max_array": max_array}).write_csv(
        f"{Path(__file__).parents[3]}/ml_service/preprocess/scaler_tensors.csv"
    )

    autoencoder = AutoEncoder(
        input_dim=52,
        hidden_dim=24,
        min_array=min_array,
        max_array=max_array,
    ).to(device)
    x_train = autoencoder.min_max_scaling(torch.tensor(x_train, dtype=torch.float32))

    criterion = nn.MSELoss()
    optimizer = optim.Adam(autoencoder.parameters(), lr=0.001)

    train_dataset = TensorDataset(x_train)
    train_loader = DataLoader(train_dataset, batch_size=50, shuffle=True)

    training_loss_list = list()

    # Early stopping variables
    best_loss = float("inf")
    patience_counter = 0

    for epoch in range(epochs):
        training_loss = 0.0
        for data in train_loader:
            input_data = data[0].to(device)
            optimizer.zero_grad()
            outputs = autoencoder(input_data)
            loss = criterion(outputs, input_data)
            loss.backward()
            optimizer.step()
            training_loss += loss.item()

        average_training_loss = training_loss / len(train_loader)
        training_loss_list.append(average_training_loss)

        if (epoch + 1) % 2 == 0:
            logger.info(f"Epoch: {epoch + 1}, train loss: {average_training_loss:.4f}")

        # Early stopping logic
        if average_training_loss < best_loss:
            best_loss = average_training_loss
            patience_counter = 0
        else:
            patience_counter += 1

        if patience_counter >= patience:
            logger.info(f"Early stopping triggered at epoch {epoch + 1}")
            break

    graph_path = "./output/training_loss.png"
    remove_graph(graph_path)
    plt.plot(range(len(training_loss_list)), training_loss_list)
    plt.savefig(graph_path)
    plt.close()

    # Save the model in onnx format
    dummy_input = torch.randn(1, 52)
    model_path = f"{Path(__file__).parents[3]}/ml_service/models/model_{version}.onnx"
    torch.onnx.export(
        autoencoder,
        dummy_input,
        model_path,
        input_names=["input"],
        output_names=["output"],
        dynamic_axes={"input": {0: "batch_size"}, "output": {0: "batch_size"}},
    )
    logger.info(f"Model saved at {model_path}")
    return autoencoder
