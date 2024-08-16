import click

from loguru import logger

from ml_model import (
    load_data,
    train_model,
    clean_data,
    filter_data_train,
    preprocess_data,
    evaluate_model,
)


@click.command
@click.option("--epochs", default=10, help="Number of epochs")
@click.option("--version", default="0.0.1", help="Model version")
def main(epochs, version):
    logger.info(f"## Training model for {epochs} epochs ##")

    logger.info("## Data Cleaning & Preprocessing ##")
    df_train = load_data("sensor.csv")

    df_train = clean_data(data=df_train)
    df_train = filter_data_train(df_train)
    x_train, _ = preprocess_data(data=df_train)

    # Train model
    logger.info("## Model Training ##")
    model = train_model(x_train=x_train, epochs=epochs, version=version)

    df_test = load_data("sensor.csv")

    logger.info("## Model Evaluation ##")
    df_eval = clean_data(data=df_test)
    x_test, _ = preprocess_data(data=df_eval)

    # Evaluate model
    _ = evaluate_model(
        df_eval=df_eval,
        x_test=x_test,
        model=model,
    )


if __name__ == "__main__":
    main()
