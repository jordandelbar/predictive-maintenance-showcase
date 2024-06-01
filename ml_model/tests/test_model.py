import numpy as np
import torch

from ml_model.model import return_normal_events_sets, AutoEncoder


def test_normal_events_sets_shape(preprocessed_test_arrays):
    # Arrange
    x, y = preprocessed_test_arrays

    # Act
    x, _ = return_normal_events_sets(x_train=x, y_train=y)

    # Assert
    assert x.shape == (2, 5)


def test_min_max_scaling(preprocessed_test_arrays):
    # Arrange
    x, _ = preprocessed_test_arrays
    min_tensor = np.array([293, 300, 20, 15, 10])
    max_tensor = np.array([297, 305, 25, 18, 12])
    autoencoder = AutoEncoder(
        input_dim=2,
        hidden_dim=2,
        min_tensor=min_tensor,
        max_tensor=max_tensor,
        min_scaling_range=0,
        max_scaling_range=1,
    )
    expected_tensor = torch.tensor(
        [
            [0.0, 0.0, 0.0, 0.0, 0.0],
            [0.0, 0.0, 0.0, 0.0, 0.0],
            [1.0, 1.0, 1.0, 1.0, 1.0],
        ]
    )

    # Act
    x_scaled = autoencoder.min_max_scaling(x=torch.tensor(x, dtype=torch.float32))

    # Assert
    assert torch.equal(x_scaled, expected_tensor)
