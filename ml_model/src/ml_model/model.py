import torch
import torch.nn as nn
import numpy as np
import polars as pl
from torch.utils.data import DataLoader
from sklearn.metrics import fbeta_score
from loguru import logger


class AutoEncoder(nn.Module):
    def __init__(
        self,
        input_dim,
        hidden_dim,
        min_tensor,
        max_tensor,
        min_scaling_range,
        max_scaling_range,
    ):
        super(AutoEncoder, self).__init__()
        self.min_tensor = torch.tensor(min_tensor, dtype=torch.float32)
        self.max_tensor = torch.tensor(max_tensor, dtype=torch.float32)
        self.min_scaling_range = min_scaling_range
        self.max_scaling_range = max_scaling_range
        self.first_encoder = nn.Sequential(
            nn.Linear(input_dim, hidden_dim),
            nn.ReLU(),
            nn.Linear(hidden_dim, hidden_dim // 2),
            nn.ReLU(),
        )
        self.first_decoder = nn.Sequential(
            nn.Linear(hidden_dim // 2, hidden_dim),
            nn.ReLU(),
            nn.Linear(hidden_dim, input_dim),
        )
        self.second_encoder = nn.Sequential(
            nn.Linear(input_dim * 2, hidden_dim * 2),
            nn.ReLU(),
            nn.Linear(hidden_dim * 2, hidden_dim),
            nn.ReLU(),
        )
        self.second_decoder = nn.Sequential(
            nn.Linear(hidden_dim, hidden_dim * 2),
            nn.ReLU(),
            nn.Linear(hidden_dim * 2, input_dim * 2),
        )
        self.third_encoder = nn.Sequential(
            nn.Linear(input_dim * 3, hidden_dim * 3),
            nn.ReLU(),
            nn.Linear(hidden_dim * 3, hidden_dim * 2),
            nn.ReLU(),
        )
        self.third_decoder = nn.Sequential(
            nn.Linear(hidden_dim * 2, hidden_dim * 3),
            nn.ReLU(),
            nn.Linear(hidden_dim * 3, input_dim),
        )

    def min_max_scaling(self, x):
        x_std = (x - self.min_tensor) / (self.max_tensor - self.min_tensor)
        x_scaled = (
            x_std * (self.max_scaling_range - self.min_scaling_range)
            + self.min_scaling_range
        )
        return x_scaled

    def forward(self, x):
        x1 = self.first_encoder(x)
        x1 = self.first_decoder(x1)
        x2 = self.second_encoder(torch.cat((x, x1), dim=1))
        x2 = self.second_decoder(x2)
        x3 = self.third_encoder(torch.cat((x, x2), dim=1))
        x3 = self.third_decoder(x3)
        return x3

    def predict(self, input_data):
        input_data = self.min_max_scaling(input_data)
        with torch.no_grad():
            output = self(input_data)
            error = torch.mean(torch.square(output - input_data), dim=1)
            return error


def train(
    x_ne_scaled,
    optimizer,
    autoencoder,
    criterion,
    device,
    epochs: int = 25,
):
    train_loader = DataLoader(x_ne_scaled, batch_size=50, shuffle=True)

    training_loss_list = list()
    for epoch in range(epochs):
        training_loss = 0.0
        for data in train_loader:
            input_data = data.to(device)
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

    return autoencoder


def return_normal_events_sets(
    x_train: np.ndarray, y_train: np.ndarray
) -> tuple[np.ndarray, np.ndarray]:
    temp_array = (
        pl.DataFrame(x_train)
        .with_columns(pl.Series(name="label", values=y_train))
        .filter(pl.col("label") == 0)
        .to_numpy()
    )

    return temp_array[:, :-1], temp_array[:, -1]


def compute_best_threshold(
    x_train: np.ndarray,
    y_train: np.ndarray,
    model: AutoEncoder,
    f_beta_threshold: float = 1.5,
) -> float:
    score_list = list()
    x_train = torch.tensor(x_train, dtype=torch.float32)
    predictions = model.predict(x_train)
    eval_df = pl.DataFrame(
        {
            "reconstruction_errors": np.array(predictions),
            "anomaly": y_train,
        }
    )

    linspace_min = np.min(eval_df["reconstruction_errors"].to_numpy())
    linspace_max = np.max(eval_df["reconstruction_errors"].to_numpy())

    for i in np.linspace(linspace_min, linspace_max, num=4000, endpoint=True):
        eval_df = eval_df.with_columns(
            pl.when(pl.col("reconstruction_errors") >= i)
            .then(1)
            .otherwise(0)
            .alias("test")
        )
        score_list.append(
            (
                i,
                fbeta_score(
                    y_true=eval_df["anomaly"],
                    y_pred=eval_df["test"],
                    beta=f_beta_threshold,
                ),
            )
        )
    threshold = sorted(score_list, key=lambda x: x[1], reverse=True)[0][0]
    return threshold
