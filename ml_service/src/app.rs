use axum::{
    routing::{get, post},
    Router,
};
use csv::Reader;
use ort::Session;
use std::error::Error;
use std::fs::File;
use std::io::BufReader;
use std::sync::{atomic::AtomicUsize, Arc};

use crate::configuration::Settings;
use crate::routes::{healthcheck, predict};

#[derive(Clone)]
pub struct AppState {
    pub sessions: Arc<Vec<Arc<Session>>>,
    pub counter: Arc<AtomicUsize>,
    pub min_values: Vec<f32>,
    pub max_values: Vec<f32>,
}

pub fn create_app(_cfg: Settings) -> Result<Router, Box<dyn Error>> {
    let num_sessions = 16;
    let sessions = (0..num_sessions)
        .map(|_| {
            let session = Session::builder()?.commit_from_file("./models/model_0.0.1.onnx")?;
            Ok(Arc::new(session))
        })
        .collect::<Result<Vec<_>, ort::Error>>()?;

    tracing::info!("created {} ONNX sessions", num_sessions);

    let (min_values, max_values) = load_scaler_tensors("./preprocess/scaler_tensors.csv")?;

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

fn load_scaler_tensors(csv_file_path: &str) -> Result<(Vec<f32>, Vec<f32>), Box<dyn Error>> {
    let file = File::open(csv_file_path)?;
    let mut rdr = Reader::from_reader(BufReader::new(file));

    let mut min_values = Vec::new();
    let mut max_values = Vec::new();

    for result in rdr.records() {
        let record = result?;
        min_values.push(record[0].parse::<f32>()?);
        max_values.push(record[1].parse::<f32>()?);
    }

    Ok((min_values, max_values))
}
