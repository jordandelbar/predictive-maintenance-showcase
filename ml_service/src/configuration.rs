use serde::Deserialize;
use std::path::PathBuf;

#[derive(Deserialize, Clone)]
pub struct Settings {
    pub service: ServiceSettings,
    pub model: ModelSettings,
}

#[derive(Deserialize, Clone)]
pub struct ServiceSettings {
    pub host: String,
    pub port: u16,
}

#[derive(Deserialize, Clone)]
pub struct ModelSettings {
    pub onnx_file: String,
    pub scaling_values_file: String,
    #[serde(default = "default_model_instances")]
    pub num_instances: usize,
    pub model_dir: PathBuf,
    pub preprocess_dir: PathBuf,
}

fn default_model_instances() -> usize {
    std::thread::available_parallelism()
        .map(|n| n.get())
        .unwrap_or(4)
}

impl ModelSettings {
    pub fn get_model_path(&self) -> PathBuf {
        self.model_dir.join(&self.onnx_file)
    }

    pub fn get_scaling_path(&self) -> PathBuf {
        self.preprocess_dir.join(&self.scaling_values_file)
    }

    pub fn validate(&self) -> Result<(), String> {
        if !self.get_model_path().exists() {
            return Err(format!("Model file not found: {:?}", self.get_model_path()));
        }
        if !self.get_scaling_path().exists() {
            return Err(format!(
                "Scaling file not found: {:?}",
                self.get_scaling_path()
            ));
        }
        Ok(())
    }
}

pub fn get_configuration() -> Result<Settings, config::ConfigError> {
    let base_path = std::env::current_dir().expect("Failed to determine the current directory");
    let configuration_directory = base_path.join("configuration");

    let environment: Environment = std::env::var("APP_ENVIRONMENT")
        .unwrap_or_else(|_| "local".into())
        .try_into()
        .expect("Failed to parse APP_ENVIRONMENT.");
    let settings = config::Config::builder()
        .add_source(config::File::from(
            configuration_directory.join("base.yaml"),
        ))
        .add_source(config::File::from(
            configuration_directory.join(format!("{}.yaml", environment.as_str())),
        ))
        .add_source(
            config::Environment::with_prefix("APP")
                .prefix_separator("_")
                .separator("__"),
        )
        .build()?;

    let settings = settings.try_deserialize::<Settings>()?;
    if let Err(e) = settings.model.validate() {
        tracing::error!("Configuration validation failed: {}", e);
        return Err(config::ConfigError::Message(e));
    }

    Ok(settings)
}

pub enum Environment {
    Local,
    Production,
}

impl Environment {
    pub fn as_str(&self) -> &'static str {
        match self {
            Environment::Local => "local",
            Environment::Production => "production",
        }
    }
}

impl TryFrom<String> for Environment {
    type Error = String;

    fn try_from(s: String) -> Result<Self, Self::Error> {
        match s.to_lowercase().as_str() {
            "local" => Ok(Self::Local),
            "production" => Ok(Self::Production),
            other => Err(format!(
                "{} is not a supported environment. Use either `local` or `production`.",
                other
            )),
        }
    }
}
