from pathlib import Path
from typing import Dict, Type, Any

import polars as pl

from loguru import logger
from polars import Int64, String, Float64


def load_data(filename: str, exclude_index: bool = True) -> pl.DataFrame:
    schema: Dict[str | Any, Type[Int64 | String | Float64] | Any] = {
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
    data = _load_csv(filename, schema=schema)

    # Exclude 'index' and casting timestamp to datetime
    if exclude_index:
        return data.select(pl.exclude("index")).with_columns(
            pl.col("timestamp").str.to_datetime().alias("timestamp")
        )
    return data.with_columns(pl.col("timestamp").str.to_datetime().alias("timestamp"))


def _load_csv(
    filename: str,
    schema=Type[Dict[str | Any, Type[pl.Int64 | pl.String | pl.Float64] | Any]],
) -> pl.DataFrame:
    logger.info(f"Loading data from {filename}")
    data = pl.read_csv(f"{Path(__file__).parents[3]}/data/{filename}", schema=schema)
    return data
