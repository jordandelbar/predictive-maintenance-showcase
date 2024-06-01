from typing import Tuple

import numpy as np
import polars as pl
import pytest


@pytest.fixture
def test_dataframe() -> pl.DataFrame:
    df = pl.DataFrame(
        {
            "Air temperature [K]": [293, 293, 297],
            "Process temperature [K]": [300, 300, 305],
            "Rotational speed [rpm]": [20, 20, 25],
            "Torque [Nm]": [15, 15, 18],
            "Tool wear [min]": [10, 10, 12],
            "Machine failure": [0, 0, 1],
        }
    )
    return df


@pytest.fixture
def preprocessed_test_arrays() -> Tuple[np.ndarray, np.ndarray]:
    x_df = pl.DataFrame(
        {
            "Air temperature [K]": [293, 293, 297],
            "Process temperature [K]": [300, 300, 305],
            "Rotational speed [rpm]": [20, 20, 25],
            "Torque [Nm]": [15, 15, 18],
            "Tool wear [min]": [10, 10, 12],
        }
    )
    y_df = pl.DataFrame(
        {
            "Machine failure": [0, 0, 1],
        }
    )
    return x_df.to_numpy(), y_df.to_numpy().squeeze()
