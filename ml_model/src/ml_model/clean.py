from typing import Dict

import polars as pl

from loguru import logger


def clean_data(data: pl.DataFrame) -> pl.DataFrame:
    # Fill null with median data
    data = _fill_nulls_with_median(data)
    if data.null_count().sum_horizontal()[0]:
        raise ValueError("There are still null values in the data")

    return data


def filter_data_train(data: pl.DataFrame) -> pl.DataFrame:
    # Get timestamps to exclude
    timestamps_to_exclude = _get_timestamps_to_exclude(data)

    # Remove data around broken status and keep NORMAL data
    _filter_out_data_train(data, timestamps_to_exclude)

    return data


def _fill_nulls_with_median(df: pl.DataFrame) -> pl.DataFrame:
    logger.info("Filling nulls with median")
    for col in df.columns:
        if df[col].dtype != pl.Utf8 and df[col].dtype != pl.Datetime:
            median_value = df[col].median()
            if median_value is None:
                median_value = 0
            df = df.with_columns(pl.col(col).fill_null(median_value))
        elif df[col].dtype == pl.Utf8:
            df = df.with_columns(pl.col(col).fill_null("MISSING"))
    return df


def _get_timestamps_to_exclude(df: pl.DataFrame) -> Dict:
    timestamps_to_exclude = (
        df.filter(pl.col("machine_status") == "BROKEN")
        .with_columns(timestamp_minus_1=pl.col("timestamp") - pl.duration(hours=1))
        .select(["timestamp", "timestamp_minus_1"])
    ).to_dict(as_series=False)

    return timestamps_to_exclude


def _filter_out_data_train(df: pl.DataFrame, timestamps_to_exclude: Dict):
    for start, end in zip(
        timestamps_to_exclude["timestamp_minus_1"], timestamps_to_exclude["timestamp"]
    ):
        # We take the data that is NOT between timestamp_minus_1 and timestamp
        df = df.filter(~((pl.col("timestamp") >= start) & (pl.col("timestamp") <= end)))
    df = df.filter(pl.col("machine_status") == "NORMAL")

    return df
