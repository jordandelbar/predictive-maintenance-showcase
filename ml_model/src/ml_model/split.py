from typing import Tuple

import numpy as np

from loguru import logger
from sklearn.model_selection import train_test_split


def split_data(
    x: np.ndarray,
    y: np.ndarray,
    test_size: float = 0.3,
    random_state: int = 42,
) -> Tuple[np.ndarray, np.ndarray, np.ndarray, np.ndarray]:
    x_train, x_test, y_train, y_test = train_test_split(
        x, y, test_size=test_size, stratify=y, random_state=random_state
    )
    logger.info(f"x train shape: {x_train.shape}")
    logger.info(f"x test shape: {x_test.shape}")
    logger.info(f"y train shape: {y_train.shape}")
    logger.info(f"y test shape: {y_test.shape}")
    return x_train, x_test, y_train, y_test
