from ml_model.preprocess import load_data, clean_data, preprocess_data
from ml_model.train import split_data, train_model, evaluate_model


def main():
    df = load_data("sensor")
    df = clean_data(data=df)
    x, y = preprocess_data(data=df)
    x_train, x_test, y_train, y_test = split_data(x=x, y=y, test_size=0.3)
    model = train_model(x_train=x_train, y_train=y_train, f_beta_score=1.4, epochs=50)
    evaluate_model(x_test=x_test, y_test=y_test, model=model)


if __name__ == "__main__":
    main()
