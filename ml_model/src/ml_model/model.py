import torch
import torch.nn as nn


class AutoEncoder(nn.Module):
    def __init__(
        self,
        input_dim,
        hidden_dim,
        min_array,
        max_array,
    ):
        super(AutoEncoder, self).__init__()
        self.min_tensor = torch.tensor(min_array, dtype=torch.float32)
        self.max_tensor = torch.tensor(max_array, dtype=torch.float32)
        self.encoder = nn.Sequential(
            nn.Linear(input_dim, hidden_dim),
            nn.ReLU(),
            nn.Linear(hidden_dim, hidden_dim // 2),
            nn.ReLU(),
        )
        self.decoder = nn.Sequential(
            nn.Linear(hidden_dim // 2, hidden_dim),
            nn.ReLU(),
            nn.Linear(hidden_dim, input_dim),
        )

    def min_max_scaling(self, x):
        x_std = (x - self.min_tensor) / (self.max_tensor - self.min_tensor)

        # In case self.min_tensor and self.max_tensor contain same vector
        x_scaled = torch.nan_to_num(x_std, nan=0.0)
        return x_scaled

    def forward(self, x):
        x = self.encoder(x)
        x = self.decoder(x)
        return x
