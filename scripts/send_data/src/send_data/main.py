import requests
import csv
import time
from pathlib import Path


# Function to send data for prediction
def send_prediction(data):
    url = "http://localhost:4000/v1/predict"
    headers = {"Content-Type": "application/json"}
    start_time = time.time()
    response = requests.post(url, headers=headers, json=data)
    end_time = time.time()
    elapsed_time = (end_time - start_time) * 1000
    print(response.text, f"{elapsed_time:.2f} ms")


def to_float(value):
    try:
        return float(value)
    except ValueError:
        return 0.0


# Read the dataset from CSV file
divergence_temperature = 0.0
divergence_rotational_speed = 0.0
divergence_torque = 0.0
count = 0
while True:
    with open(
        f"{Path(__file__).parents[4]}/data/sensor.csv",
        newline="",
        encoding="utf-8-sig",
    ) as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            # Prepare data for prediction
            data = {
                "machine_id": 7,
                "sensor_00": to_float(row["sensor_00"]),
                "sensor_01": to_float(row["sensor_01"]),
                "sensor_02": to_float(row["sensor_02"]),
                "sensor_03": to_float(row["sensor_03"]),
                "sensor_04": to_float(row["sensor_04"]),
                "sensor_05": to_float(row["sensor_05"]),
                "sensor_06": to_float(row["sensor_06"]),
                "sensor_07": to_float(row["sensor_07"]),
                "sensor_08": to_float(row["sensor_08"]),
                "sensor_09": to_float(row["sensor_09"]),
                "sensor_10": to_float(row["sensor_10"]),
                "sensor_11": to_float(row["sensor_11"]),
                "sensor_12": to_float(row["sensor_12"]),
                "sensor_13": to_float(row["sensor_13"]),
                "sensor_14": to_float(row["sensor_14"]),
                "sensor_15": to_float(row["sensor_15"]),
                "sensor_16": to_float(row["sensor_16"]),
                "sensor_17": to_float(row["sensor_17"]),
                "sensor_18": to_float(row["sensor_18"]),
                "sensor_19": to_float(row["sensor_19"]),
                "sensor_20": to_float(row["sensor_20"]),
                "sensor_21": to_float(row["sensor_21"]),
                "sensor_22": to_float(row["sensor_22"]),
                "sensor_23": to_float(row["sensor_23"]),
                "sensor_24": to_float(row["sensor_24"]),
                "sensor_25": to_float(row["sensor_25"]),
                "sensor_26": to_float(row["sensor_26"]),
                "sensor_27": to_float(row["sensor_27"]),
                "sensor_28": to_float(row["sensor_28"]),
                "sensor_29": to_float(row["sensor_29"]),
                "sensor_30": to_float(row["sensor_30"]),
                "sensor_31": to_float(row["sensor_31"]),
                "sensor_32": to_float(row["sensor_32"]),
                "sensor_33": to_float(row["sensor_33"]),
                "sensor_34": to_float(row["sensor_34"]),
                "sensor_35": to_float(row["sensor_35"]),
                "sensor_36": to_float(row["sensor_36"]),
                "sensor_37": to_float(row["sensor_37"]),
                "sensor_38": to_float(row["sensor_38"]),
                "sensor_39": to_float(row["sensor_39"]),
                "sensor_40": to_float(row["sensor_40"]),
                "sensor_41": to_float(row["sensor_41"]),
                "sensor_42": to_float(row["sensor_42"]),
                "sensor_43": to_float(row["sensor_43"]),
                "sensor_44": to_float(row["sensor_44"]),
                "sensor_45": to_float(row["sensor_45"]),
                "sensor_46": to_float(row["sensor_46"]),
                "sensor_47": to_float(row["sensor_47"]),
                "sensor_48": to_float(row["sensor_48"]),
                "sensor_49": to_float(row["sensor_49"]),
                "sensor_50": to_float(row["sensor_50"]),
                "sensor_51": to_float(row["sensor_51"]),
            }
            # Send data for prediction
            send_prediction(data)
