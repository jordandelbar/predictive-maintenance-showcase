from pathlib import Path

from ml_model.preprocess import (
    load_data,
    clean_data_train,
    clean_data_evaluate,
    preprocess_data,
)
from ml_model.train import train_model, evaluate_model


def main():
    df = load_data("sensor")
    df_train = clean_data_train(data=df, lower_quantile=0.01, upper_quantile=0.99)
    x, y = preprocess_data(data=df_train)
    model = train_model(x_train=x, y_train=y, epochs=100)
    df_eval = clean_data_evaluate(data=df)
    x, y = preprocess_data(data=df_eval)
    reconstruction_error = evaluate_model(x_test=x, model=model)
    df_eval.with_columns(reconstruction_error=reconstruction_error).write_csv(
        f"{Path(__file__).parents[2]}/output/reconstruction_error.csv"
    )


if __name__ == "__main__":
    main()
