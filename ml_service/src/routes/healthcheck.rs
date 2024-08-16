use axum::response::{IntoResponse, Json};
use serde::{Deserialize, Serialize};
use tracing::instrument;

#[derive(Serialize, Deserialize)]
pub struct Status {
    status: String,
}

#[instrument()]
pub async fn healthcheck() -> impl IntoResponse {
    Json(Status {
        status: "Available".into(),
    })
}
