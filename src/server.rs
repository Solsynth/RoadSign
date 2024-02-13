use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct ServerBindConfig {
    pub addr: String,
    pub tls: bool,
}