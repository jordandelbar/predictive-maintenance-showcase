import numpy as np
import torch

from ml_model import AutoEncoder


def test_min_max_scaling(preprocessed_test_arrays):
    # Arrange
    x, _ = preprocessed_test_arrays

    min_array = np.array([1.0 * i for i in range(1, 53)])
    max_array = np.array([1.2 * i for i in range(1, 53)])

    autoencoder = AutoEncoder(
        input_dim=2,
        hidden_dim=2,
        min_array=min_array,
        max_array=max_array,
    )
    expected_tensor = torch.tensor(
        [
            [0.0000] * 52,
            [0.5000] * 52,
            [1.0000] * 52,
        ]
    )

    # Act
    x_scaled = autoencoder.min_max_scaling(
        x=torch.tensor(x, dtype=torch.float32),
    )

    # Assert
    assert torch.all(torch.isclose(x_scaled, expected_tensor))


def test_min_max_scaling_same_data(preprocessed_test_arrays_same):
    # Arrange
    x, _ = preprocessed_test_arrays_same

    min_array = np.array([1.0 * i for i in range(1, 53)])
    max_array = np.array([1.0 * i for i in range(1, 53)])

    autoencoder = AutoEncoder(
        input_dim=2,
        hidden_dim=2,
        min_array=min_array,
        max_array=max_array,
    )
    expected_tensor = torch.tensor(
        [
            [0.0000] * 52,
            [0.0000] * 52,
            [0.0000] * 52,
        ]
    )

    # Act
    x_scaled = autoencoder.min_max_scaling(
        x=torch.tensor(x, dtype=torch.float32),
    )

    # Assert
    assert torch.all(torch.isclose(x_scaled, expected_tensor))
