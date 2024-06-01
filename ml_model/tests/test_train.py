import numpy as np

from ml_model.train import preprocess_data


def test_preprocess_data_correct_shape(test_dataframe):
    # Arrange
    df = test_dataframe

    # Act
    x, y = preprocess_data(df)

    # Assert
    assert x.shape == (3, 5)
    assert y.shape == (3,)


def test_preprocess_data_correct_labels(test_dataframe):
    # Arrange
    df = test_dataframe

    # Act
    _, y = preprocess_data(df)

    # Assert
    assert np.sum(y) == 1


def test_preprocess_data_equal_arrays(test_dataframe, preprocessed_test_arrays):
    # Arrange
    df = test_dataframe
    x_expected, y_expected = preprocessed_test_arrays

    # Act
    x, y = preprocess_data(df)

    # Assert
    np.testing.assert_array_equal(x, x_expected)
    np.testing.assert_array_equal(y, y_expected)
