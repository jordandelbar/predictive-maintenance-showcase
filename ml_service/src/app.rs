use crate::configuration::Settings;
use crate::routes::{healthcheck, predict};
use axum::{
    routing::{get, post},
    Router,
};
use csv::Reader;
use ndarray::{Array, Array1};
use ort::Session;
use std::error::Error;
use std::fs::File;
use std::io::BufReader;
use std::path::Path;
use std::sync::{atomic::AtomicUsize, Arc};

#[derive(Clone)]
pub struct AppState {
    pub sessions: Arc<Vec<Arc<Session>>>,
    pub counter: Arc<AtomicUsize>,
    pub min_values: Array1<f32>,
    pub max_values: Array1<f32>,
}

pub fn create_app(cfg: Settings) -> Result<Router, Box<dyn Error>> {
    let sessions = (0..cfg.model.num_instances)
        .map(|_| {
            let session = Session::builder()?.commit_from_file(cfg.model.get_model_path())?;
            Ok(Arc::new(session))
        })
        .collect::<Result<Vec<_>, ort::Error>>()?;

    tracing::info!("created {} ONNX sessions", cfg.model.num_instances);

    let (min_values, max_values) = load_scaler_tensors(&cfg.model.get_scaling_path())?;

    let app_state = AppState {
        sessions: Arc::new(sessions),
        counter: Arc::new(AtomicUsize::new(0)),
        min_values,
        max_values,
    };

    let app = Router::new()
        .route("/health", get(healthcheck))
        .route("/predict", post(predict))
        .with_state(app_state);
    Ok(app)
}

fn load_scaler_tensors(csv_file_path: &Path) -> Result<(Array1<f32>, Array1<f32>), Box<dyn Error>> {
    let file = File::open(csv_file_path).map_err(|e| {
        Box::new(std::io::Error::new(
            std::io::ErrorKind::Other,
            format!("Failed to open scaling file {:?}: {}", csv_file_path, e),
        ))
    })?;

    let mut rdr = Reader::from_reader(BufReader::new(file));
    let mut min_values = Vec::new();
    let mut max_values = Vec::new();

    for (idx, result) in rdr.records().enumerate() {
        let record = result.map_err(|e| {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::Other,
                format!("Failed to read record at line {}: {}", idx + 1, e),
            ))
        })?;

        if record.len() != 2 {
            return Err(Box::new(std::io::Error::new(
                std::io::ErrorKind::InvalidData,
                format!(
                    "Invalid record at line {}: expected 2 values, got {}",
                    idx + 1,
                    record.len()
                ),
            )));
        }

        min_values.push(record[0].parse::<f32>().map_err(|e| {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::InvalidData,
                format!("Failed to parse min value at line {}: {}", idx + 1, e),
            ))
        })?);

        max_values.push(record[1].parse::<f32>().map_err(|e| {
            Box::new(std::io::Error::new(
                std::io::ErrorKind::InvalidData,
                format!("Failed to parse max value at line {}: {}", idx + 1, e),
            ))
        })?);
    }

    if min_values.is_empty() {
        return Err(Box::new(std::io::Error::new(
            std::io::ErrorKind::InvalidData,
            "Scaling file is empty",
        )));
    }

    Ok((Array::from_vec(min_values), Array::from_vec(max_values)))
}
