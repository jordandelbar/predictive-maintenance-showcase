from typing import Tuple

import numpy as np
import polars as pl

from loguru import logger


def preprocess_data(data: pl.DataFrame) -> Tuple[np.ndarray, np.ndarray]:
    training_data = data.select(
        [f"sensor_{str(x).zfill(2)}" for x in range(52)] + ["machine_status"]
    )
    x_df = training_data.select(pl.exclude("machine_status"))
    y_df = training_data.select(pl.col("machine_status"))

    logger.info(f"x dataframe shape: {x_df.shape}")
    logger.info(f"y dataframe shape: {y_df.shape}")
    x_array = x_df.to_numpy()
    y_array = y_df.to_numpy().squeeze()

    return x_array, y_array
