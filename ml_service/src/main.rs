use ml_service::app::create_app;
use ml_service::configuration::get_configuration;
use std::error::Error;
use tokio::{net::TcpListener, signal};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "info,ort=info".into()),
        )
        .with(
            tracing_subscriber::fmt::layer()
                .json()
                .with_target(false)
                .with_level(true)
                .with_thread_names(true)
                .with_thread_ids(true),
        )
        .init();

    let cfg = get_configuration().expect("Failed to parse configuration");

    tracing::info!("starting the ml service app");

    let addr = format!("{}:{}", cfg.service.host, cfg.service.port);
    let app = create_app(cfg)
        .expect("failed to create app")
        .into_make_service();
    let listener = TcpListener::bind(&addr).await?;

    tracing::info!("listening on {}", &addr);
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal())
        .await?;

    Ok(())
}

async fn shutdown_signal() {
    signal::ctrl_c()
        .await
        .expect("Failed to install CTRL+C signal handler");
    tracing::info!("Received CTRL+C signal, shutting down...");
}
