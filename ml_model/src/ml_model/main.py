from train import load_csv, preprocess_data, split_data, train_model, evaluate_model


def main():
    df = load_csv("ai4i2020")
    x, y = preprocess_data(data=df)
    x_train, x_test, y_train, y_test = split_data(x=x, y=y, test_size=0.3)
    model, best_threshold = train_model(
        x_train=x_train, y_train=y_train, f_beta_score=1.4, epochs=100
    )
    evaluate_model(x_test=x_test, y_test=y_test, model=model, threshold=best_threshold)


if __name__ == "__main__":
    main()
