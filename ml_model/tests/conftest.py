from typing import Tuple

import numpy as np
import polars as pl
import pytest


@pytest.fixture()
def data():
    return {
        f"sensor_{str(i).zfill(2)}": [1.0 * j, 1.1 * j, 1.2 * j]
        for i, j in zip(range(52), range(1, 53))
    }


@pytest.fixture()
def data_same():
    return {
        f"sensor_{str(i).zfill(2)}": [1.0 * j, 1.0 * j, 1.0 * j]
        for i, j in zip(range(52), range(1, 53))
    }


@pytest.fixture
def test_dataframe(data) -> pl.DataFrame:
    df = pl.DataFrame(data).with_columns(
        pl.Series(name="machine_status", values=[0, 0, 1])
    )
    print(df)
    return df


@pytest.fixture
def preprocessed_test_arrays(data) -> Tuple[np.ndarray, np.ndarray]:
    x_df = pl.DataFrame(data)
    y_df = pl.DataFrame(
        {
            "machine_status": [0, 0, 1],
        }
    )
    return x_df.to_numpy(), y_df.to_numpy().squeeze()


@pytest.fixture
def preprocessed_test_arrays_same(data_same) -> Tuple[np.ndarray, np.ndarray]:
    x_df = pl.DataFrame(data_same)
    y_df = pl.DataFrame(
        {
            "machine_status": [0, 0, 1],
        }
    )
    return x_df.to_numpy(), y_df.to_numpy().squeeze()
