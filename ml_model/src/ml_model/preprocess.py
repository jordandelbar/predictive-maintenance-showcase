from typing import Tuple

import numpy as np
import polars as pl

from loguru import logger

from ml_model.helpers import (
    load_csv,
    fill_nulls_with_median,
    get_quantiles,
    filter_by_quantiles,
)


def load_data(filename: str) -> pl.DataFrame:
    schema = {
        "index": pl.Int64,
        "timestamp": pl.Utf8,
        "sensor_00": pl.Float64,
        "sensor_01": pl.Float64,
        "sensor_02": pl.Float64,
        "sensor_03": pl.Float64,
        "sensor_04": pl.Float64,
        "sensor_05": pl.Float64,
        "sensor_06": pl.Float64,
        "sensor_07": pl.Float64,
        "sensor_08": pl.Float64,
        "sensor_09": pl.Float64,
        "sensor_10": pl.Float64,
        "sensor_11": pl.Float64,
        "sensor_12": pl.Float64,
        "sensor_13": pl.Float64,
        "sensor_14": pl.Float64,
        "sensor_15": pl.Float64,
        "sensor_16": pl.Float64,
        "sensor_17": pl.Float64,
        "sensor_18": pl.Float64,
        "sensor_19": pl.Float64,
        "sensor_20": pl.Float64,
        "sensor_21": pl.Float64,
        "sensor_22": pl.Float64,
        "sensor_23": pl.Float64,
        "sensor_24": pl.Float64,
        "sensor_25": pl.Float64,
        "sensor_26": pl.Float64,
        "sensor_27": pl.Float64,
        "sensor_28": pl.Float64,
        "sensor_29": pl.Float64,
        "sensor_30": pl.Float64,
        "sensor_31": pl.Float64,
        "sensor_32": pl.Float64,
        "sensor_33": pl.Float64,
        "sensor_34": pl.Float64,
        "sensor_35": pl.Float64,
        "sensor_36": pl.Float64,
        "sensor_37": pl.Float64,
        "sensor_38": pl.Float64,
        "sensor_39": pl.Float64,
        "sensor_40": pl.Float64,
        "sensor_41": pl.Float64,
        "sensor_42": pl.Float64,
        "sensor_43": pl.Float64,
        "sensor_44": pl.Float64,
        "sensor_45": pl.Float64,
        "sensor_46": pl.Float64,
        "sensor_47": pl.Float64,
        "sensor_48": pl.Float64,
        "sensor_49": pl.Float64,
        "sensor_50": pl.Float64,
        "sensor_51": pl.Float64,
        "machine_status": pl.String,
    }
    data = load_csv(filename, schema=schema)

    # Exclude 'index' and casting timestamp to datetime
    data = data.select(pl.exclude("index"))
    data = data.with_columns(pl.col("timestamp").str.to_datetime().alias("timestamp"))
    return data


def clean_data_train(
    data: pl.DataFrame, lower_quantile: float = 0.02, upper_quantile: float = 0.98
) -> pl.DataFrame:
    # Only taking NORMAL data
    data = data.filter(pl.col("machine_status") == "NORMAL")
    # Fill null with median data
    data = fill_nulls_with_median(data)
    if data.null_count().sum(axis=1)[0]:
        raise ValueError("There are still null values in the data")
    # Remove outliers
    quantiles = get_quantiles(data, lower_quantile, upper_quantile)
    data = filter_by_quantiles(data, quantiles)
    return data


def clean_data_evaluate(data: pl.DataFrame) -> pl.DataFrame:
    # Fill null with median data
    data = fill_nulls_with_median(data)
    if data.null_count().sum(axis=1)[0]:
        raise ValueError("There are still null values in the data")
    return data


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
