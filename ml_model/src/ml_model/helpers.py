from pathlib import Path
from typing import Dict, Tuple

import polars as pl
from loguru import logger


def load_csv(filename: str, schema=Dict[str, pl.DataType]) -> pl.DataFrame:
    logger.info(f"Loading data from {filename}")
    data = pl.read_csv(
        f"{Path(__file__).parents[3]}/data/{filename}.csv", schema=schema
    )
    return data


def fill_nulls_with_median(df: pl.DataFrame) -> pl.DataFrame:
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


def get_quantiles(
    df: pl.DataFrame, lower_bound: float, upper_bound: float
) -> Dict[str, Tuple[float, float]]:
    logger.info("Computing quantiles")
    sensor_limits = dict()
    for col in df.columns:
        if col not in ["timestamp", "machine_status", ""]:
            sensor_limits[col] = (
                df.quantile(lower_bound, "lower")[col][0],
                df.quantile(upper_bound, "higher")[col][0],
            )
    filtered_sensor_limits = {k: v for k, v in sensor_limits.items() if None not in v}
    return filtered_sensor_limits


def filter_by_quantiles(
    df: pl.DataFrame, boundaries: Dict[str, Tuple[float, float]]
) -> pl.DataFrame:
    logger.info("Filtering by quantiles")
    mask = pl.Series([True] * len(df))
    for col, (lower, upper) in boundaries.items():
        mask &= (df[col] >= lower) & (df[col] <= upper)
    return df.filter(mask)
