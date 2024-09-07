import os
from pathlib import Path
from typing import List

import matplotlib.pyplot as plt
import numpy as np
import polars as pl
import torch

from loguru import logger

from ml_model.model import AutoEncoder
from ml_model.utils import remove_graph


def evaluate_model(
    x_test: np.ndarray,
    model: AutoEncoder,
) -> List[float]:
    with torch.no_grad():
        x_test = torch.tensor(x_test, dtype=torch.float32)
        x_test = model.min_max_scaling(x_test)
        outputs = model(x_test)
        reconstruction_errors = (
            torch.mean(torch.square(outputs - x_test), dim=1).detach().numpy()
        )

    logger.info(f"Reconstruction errors: {reconstruction_errors}")

    # Plot and save output data
    _plot_breakdown_graphs()
    _plot_reconstruction_errors()

    return reconstruction_errors


def _plot_breakdown_graphs():
    logger.info("Plot breakdown graphs")
    df = (
        pl.read_csv("./output/reconstruction_errors.csv")
        .drop_nulls()
        .with_columns(pl.col("timestamp").str.to_datetime().alias("timestamp"))
    )

    # Compute hours around broken machine status
    broken_timestamps = df.filter(pl.col("machine_status") == "BROKEN")[
        "timestamp"
    ].to_list()
    broken_timestamps_series = pl.Series(broken_timestamps)
    time_windows = [
        (dt - pl.duration(hours=8), dt + pl.duration(hours=3))
        for dt in broken_timestamps_series
    ]
    window_dfs = [
        df.filter((pl.col("timestamp") >= start) & (pl.col("timestamp") <= end))
        for start, end in time_windows
    ]
    breakdown_nb = 1
    for df in window_dfs:
        normal = df.filter(pl.col("machine_status") == "NORMAL")
        broken = df.filter(pl.col("machine_status") == "BROKEN")
        recovering = df.filter(pl.col("machine_status") == "RECOVERING")
        for temp_df, style in zip([normal, broken, recovering], ["g-", "rx", "y-"]):
            plt.plot(temp_df["timestamp"], temp_df["reconstruction_errors"], style)

        graph_path = (
            f"{Path(__file__).parents[2]}/output/broken_status_graph_{breakdown_nb}.png"
        )
        remove_graph(graph_path)
        plt.savefig(graph_path)
        plt.close()

        breakdown_nb += 1


def _plot_reconstruction_errors():
    logger.info("Plot reconstruction errors graph")
    df = (
        pl.read_csv("./output/reconstruction_errors.csv")
        .drop_nulls()
        .with_columns(pl.col("timestamp").str.to_datetime().alias("timestamp"))
    )
    rec_errors = df.select(["timestamp", "reconstruction_errors", "machine_status"])
    rec_broken = rec_errors.filter(pl.col("machine_status") == "BROKEN")
    plt.plot(rec_errors["timestamp"], rec_errors["reconstruction_errors"], "g-")
    plt.plot(rec_broken["timestamp"], rec_broken["reconstruction_errors"], "rx")

    graph_path = f"{Path(__file__).parents[2]}/output/reconstruction_errors_graph.png"
    remove_graph(graph_path)
    plt.savefig(graph_path)
    plt.close()
