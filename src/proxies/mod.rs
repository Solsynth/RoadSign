use serde::{Deserialize, Serialize};

use self::config::Region;

pub mod config;
pub mod loader;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Instance {
    pub regions: Vec<Region>,
}

impl Instance {
    pub fn new() -> Instance {
        Instance { regions: vec![] }
    }
}
