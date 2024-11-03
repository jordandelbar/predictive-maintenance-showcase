from pathlib import Path

import polars as pl
import matplotlib.pyplot as plt
import seaborn as sns


def main():
    uri = "postgresql://monitor:test@localhost:5432/monitoring"
    training_df = (
        pl.read_csv(
            f"{Path(__file__).parents[2]}/ml_model/output/reconstruction_errors.csv"
        )
        .with_columns(
            pl.col("timestamp").str.strptime(
                pl.Datetime(time_unit="ns"), format="%Y-%m-%dT%H:%M:%S%.f", strict=False
            )
        )
        .rename({"reconstruction_errors": "rce_training"})
    )
    inferred_df = pl.read_database_uri("select * from monitoring", uri=uri).rename(
        {"reconstruction_error": "rce_inference"}
    )
    join_df = (
        training_df.select(["timestamp", "rce_training"])
        .join(
            inferred_df.select(["created_at", "rce_inference"]),
            how="inner",
            left_on="timestamp",
            right_on="created_at",
        )
        .with_columns(difference=pl.col("rce_training") - pl.col("rce_inference"))
        .sort(by="rce_training")
    )
    print(join_df)
    mean_difference = join_df["difference"].mean()
    median_difference = join_df["difference"].median()
    quantile_75 = join_df["difference"].quantile(0.75)
    quantile_95 = join_df["difference"].quantile(0.95)
    quantile_99 = join_df["difference"].quantile(0.99)
    quantile_999 = join_df["difference"].quantile(0.999)

    print("Mean of difference:", mean_difference)
    print("Median of difference:", median_difference)
    print("75th percentile of difference:", quantile_75)
    print("95th percentile of difference:", quantile_95)
    print("99th percentile of difference:", quantile_99)
    print("99.9th percentile of difference:", quantile_999)
    plot_differences(join_df)


def plot_differences(join_df: pl.DataFrame):
    df_pandas = join_df.to_pandas()
    plt.figure(figsize=(10, 6))
    sns.histplot(df_pandas["difference"], bins=50, kde=True)

    plt.title("Distribution of Differences")
    plt.xlabel("Difference")
    plt.ylabel("Frequency")
    plt.savefig("./distribution_of_differences.png", dpi=300)
    plt.close()


if __name__ == "__main__":
    main()
