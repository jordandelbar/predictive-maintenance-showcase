use axum::extract::State;
use axum::response::IntoResponse;
use axum::Json;
use ort::{Session, Tensor};
use serde::{Deserialize, Serialize};
use std::error::Error;
use std::sync::atomic::Ordering;
use std::sync::Arc;
use tracing::instrument;

use crate::app::AppState;

#[derive(Deserialize)]
pub struct PredictRequest {
    input_values: Vec<f32>,
}

#[derive(Deserialize, Serialize)]
pub struct PredictResponse {
    reconstruction_error: f32,
}

#[instrument(skip(state, payload))]
pub async fn predict(
    State(state): State<AppState>,
    Json(payload): Json<PredictRequest>,
) -> impl IntoResponse {
    let input_values = payload.input_values.clone();
    let scaled_input_values =
        scale_input_values(&input_values, &state.min_values, &state.max_values);

    let session_index = state.counter.fetch_add(1, Ordering::SeqCst) % state.sessions.len();
    let session = state.sessions[session_index].clone();

    tracing::info!("handling request with session {}", session_index);

    let (reconstruction_error, _output_values) = match get_outputs(session, &scaled_input_values) {
        Ok(result) => result,
        Err(err) => {
            tracing::error!(
                "error during prediction: {}, inputs: {:?}",
                err,
                &input_values
            );
            return Json(PredictResponse {
                reconstruction_error: 0.0,
            });
        }
    };

    Json(PredictResponse {
        reconstruction_error,
    })
}

fn scale_input_values(input_values: &[f32], min_values: &[f32], max_values: &[f32]) -> Vec<f32> {
    input_values
        .iter()
        .enumerate()
        .map(|(i, &x)| {
            let min_val = min_values[i];
            let max_val = max_values[i];
            let range = max_val - min_val;

            // Compute x_std
            if range != 0.0 {
                (x - min_val) / range
            } else {
                // If range is zero, return 0.0 to avoid NaN
                0.0
            }
        })
        .collect()
}

fn get_outputs(
    session: Arc<Session>,
    input_values: &[f32],
) -> Result<(f32, Vec<f32>), Box<dyn Error>> {
    let input_shape = [1usize, 52usize];
    let value = Tensor::from_array((input_shape, input_values.to_owned().into_boxed_slice()))?;

    let output_tensor = &session.run(ort::inputs![value]?)?[0];

    let output_array = output_tensor.try_extract_tensor::<f32>()?;
    let output_values = output_array
        .as_slice()
        .ok_or("Failed to convert output to slice")?
        .to_vec();

    let mse = compute_mse(input_values, &output_values)?;

    Ok((mse, output_values))
}

fn compute_mse(input: &[f32], output: &[f32]) -> Result<f32, Box<dyn Error>> {
    if input.len() != output.len() {
        return Err("Input and output arrays must have the same length".into());
    }

    let sum_of_squares: f32 = input
        .iter()
        .zip(output.iter())
        .map(|(i, o)| (i - o).powi(2)) // Square the difference
        .sum();

    let mse = sum_of_squares / input.len() as f32;

    Ok(mse)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_compute_mse_basic() {
        let input = vec![1.0, 2.0, 3.0];
        let output = vec![1.0, 2.0, 3.0];
        let expected = 0.0;
        let result = compute_mse(&input, &output).unwrap();
        assert_eq!(result, expected);
    }

    #[test]
    fn test_compute_mse_nonzero() {
        let input = vec![1.0, 2.0, 3.0];
        let output = vec![1.0, 2.0, 4.0];
        let expected = 0.33333334; // (0^2 + 0^2 + 1^2) / 3 = 0.33333334
        let result = compute_mse(&input, &output).unwrap();
        assert!((result - expected).abs() < 1e-6);
    }

    #[test]
    fn test_compute_mse_different_lengths() {
        let input = vec![1.0, 2.0, 3.0];
        let output = vec![1.0, 2.0];
        let result = compute_mse(&input, &output);
        assert!(result.is_err());
        assert_eq!(
            result.unwrap_err().to_string(),
            "Input and output arrays must have the same length"
        );
    }

    #[test]
    fn test_min_max_scaling() {
        // Arrange
        let input_values = vec![vec![1.0; 52], vec![1.1; 52], vec![1.2; 52]];

        let min_values = vec![1.0; 52];
        let max_values = vec![1.2; 52];

        // Expected output
        let expected_scaled_values = vec![vec![0.0; 52], vec![0.5; 52], vec![1.0; 52]];

        // Act & Assert
        for (input, expected) in input_values.iter().zip(expected_scaled_values.iter()) {
            let scaled_values = scale_input_values(input, &min_values, &max_values);
            assert_eq!(
                scaled_values, *expected,
                "Scaled values did not match expected values"
            );
        }
    }

    #[test]
    fn test_min_max_scaling_same_data() {
        // Arrange
        let input_values = vec![vec![1.0; 52], vec![1.0; 52], vec![1.0; 52]];

        let min_values = vec![1.0; 52];
        let max_values = vec![1.2; 52];

        // Expected output
        let expected_scaled_values = vec![vec![0.0; 52], vec![0.0; 52], vec![0.0; 52]];

        // Act & Assert
        for (input, expected) in input_values.iter().zip(expected_scaled_values.iter()) {
            let scaled_values = scale_input_values(input, &min_values, &max_values);
            assert_eq!(
                scaled_values, *expected,
                "Scaled values did not match expected values"
            );
        }
    }

    #[test]
    fn test_scale_input_values() {
        let min_values = vec![
            0.0,
            0.0,
            33.15972,
            31.640620000000002,
            2.798032,
            0.0,
            0.01446759,
            0.0,
            0.02893518,
            0.0,
            0.0,
            0.0,
            0.0,
            0.0,
            32.40955,
            0.0,
            0.0,
            0.0,
            0.0,
            0.0,
            0.0,
            95.52766,
            0.0,
            0.0,
            0.0,
            0.0,
            43.154790000000006,
            0.0,
            4.3193470000000005,
            0.6365742,
            0.0,
            23.95833,
            0.24071610000000002,
            6.460602,
            54.882369999999995,
            0.0,
            2.26097,
            0.0,
            24.4791660308838,
            19.27083,
            23.4375,
            20.83333,
            22.1354160308838,
            24.4791660308838,
            25.752315521240202,
            26.331018447876,
            26.331018447876,
            27.199070000000003,
            26.331018447876,
            26.62037,
            27.488426208496104,
            27.7777786254883,
        ];

        let max_values = vec![
            2.549016,
            56.727430000000005,
            56.032990000000005,
            48.220490000000005,
            800.0,
            99.99988,
            22.251160000000002,
            23.59664,
            24.34896,
            25.0,
            76.10686,
            60.0,
            45.0,
            31.18755,
            500.0,
            0.0,
            739.7415,
            599.999938964844,
            4.87325,
            878.9179,
            448.9079,
            1107.526,
            594.0611,
            1227.5639999999999,
            1000.0,
            839.575,
            1214.42,
            2000.0,
            1841.146,
            1466.281,
            1600.0,
            1800.0,
            1839.211,
            1578.6,
            425.5498,
            694.479125976563,
            984.0607,
            174.9012,
            417.7083,
            547.9166,
            512.7604,
            420.3125,
            374.2188,
            408.5937,
            1000.0,
            320.3125,
            370.3704,
            303.5301,
            561.632,
            464.4097,
            1000.0,
            1000.0,
        ];

        let input_values = vec![
            2.465394,
            47.092009999999995,
            53.2118,
            46.310759999999995,
            634.375,
            76.45975,
            13.41146,
            16.13136,
            15.567129999999999,
            15.053529999999999,
            37.2274,
            47.52422,
            31.11716,
            1.6813529999999999,
            419.5747,
            0.0,
            461.8781,
            466.3284,
            2.565284,
            665.3993,
            398.9862,
            880.0001,
            498.8926,
            975.9409,
            627.674,
            741.7151,
            848.0708,
            429.0377,
            785.1935,
            684.9443,
            594.4445,
            682.8125,
            680.4416,
            433.7037,
            171.9375,
            341.9039,
            195.0655,
            90.32386,
            40.36458,
            31.51042,
            70.57291,
            30.98958,
            31.770832061767603,
            41.92708,
            39.6412,
            65.68287,
            50.92593,
            38.19444,
            157.9861,
            67.70834,
            243.0556,
            201.3889,
        ];

        let expected_scaled_values = vec![
            0.9671944,
            0.83014536,
            0.8766599,
            0.8848164,
            0.7922421,
            0.7645984,
            0.60247236,
            0.6836295,
            0.63890535,
            0.6021412,
            0.48914647,
            0.7920703,
            0.69149244,
            0.05391103,
            0.82800055,
            0.0,
            0.6243777,
            0.77721405,
            0.52640104,
            0.7570665,
            0.888793,
            0.7751717,
            0.8398002,
            0.7950224,
            0.627674,
            0.8834411,
            0.68721926,
            0.21451885,
            0.4251213,
            0.46689883,
            0.37152782,
            0.37096778,
            0.36988136,
            0.27175903,
            0.31579557,
            0.49231702,
            0.19637868,
            0.51642793,
            0.040397342,
            0.023152722,
            0.09632782,
            0.02542373,
            0.027366856,
            0.045423724,
            0.014256011,
            0.13385828,
            0.071488656,
            0.03979057,
            0.24594587,
            0.0938533,
            0.22166026,
            0.17857143,
        ];

        let scaled_values = scale_input_values(&input_values, &min_values, &max_values);

        for (scaled, expected) in scaled_values.iter().zip(expected_scaled_values.iter()) {
            assert!(
                (scaled - expected).abs() < 1e-6,
                "Scaled value {} did not match expected value {}",
                scaled,
                expected
            );
        }
    }
}
