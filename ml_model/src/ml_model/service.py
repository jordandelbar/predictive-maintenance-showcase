import bentoml
import numpy as np

from bentoml.io import JSON
from pydantic import BaseModel

pm_runner = bentoml.pytorch.get("pm_autoencoder:latest").to_runner()

svc = bentoml.Service(name="pm_demo", runners=[pm_runner])


class SensorData(BaseModel):
    sensor_00: float
    sensor_01: float
    sensor_02: float
    sensor_03: float
    sensor_04: float
    sensor_05: float
    sensor_06: float
    sensor_07: float
    sensor_08: float
    sensor_09: float
    sensor_10: float
    sensor_11: float
    sensor_12: float
    sensor_13: float
    sensor_14: float
    sensor_15: float
    sensor_16: float
    sensor_17: float
    sensor_18: float
    sensor_19: float
    sensor_20: float
    sensor_21: float
    sensor_22: float
    sensor_23: float
    sensor_24: float
    sensor_25: float
    sensor_26: float
    sensor_27: float
    sensor_28: float
    sensor_29: float
    sensor_30: float
    sensor_31: float
    sensor_32: float
    sensor_33: float
    sensor_34: float
    sensor_35: float
    sensor_36: float
    sensor_37: float
    sensor_38: float
    sensor_39: float
    sensor_40: float
    sensor_41: float
    sensor_42: float
    sensor_43: float
    sensor_44: float
    sensor_45: float
    sensor_46: float
    sensor_47: float
    sensor_48: float
    sensor_49: float
    sensor_50: float
    sensor_51: float


class ModelOutput(BaseModel):
    reconstruction_error: float


def to_numpy(tensor):
    return tensor.detach().cpu().numpy()


@svc.api(input=JSON(pydantic_model=SensorData), output=JSON(pydantic_model=ModelOutput))
async def predict(inp: SensorData) -> ModelOutput:
    sensor_array = np.array([[getattr(inp, f"sensor_{i:02d}") for i in range(52)]])

    output = await pm_runner.async_run(sensor_array)
    return ModelOutput(reconstruction_error=to_numpy(output)[0])
