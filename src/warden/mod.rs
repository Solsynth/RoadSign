pub mod runner;

use std::collections::HashMap;

use futures_util::lock::Mutex;
use lazy_static::lazy_static;
use poem_openapi::Object;
use serde::{Deserialize, Serialize};
use tracing::{debug, warn};

use crate::proxies::config::Region;

use self::runner::AppInstance;

lazy_static! {
    static ref INSTANCES: Mutex<HashMap<String, AppInstance>> = Mutex::new(HashMap::new());
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WardenInstance {
    pub applications: Vec<Application>,
}

impl WardenInstance {
    pub fn new() -> WardenInstance {
        WardenInstance {
            applications: vec![],
        }
    }

    pub fn scan(&mut self, regions: Vec<Region>) {
        self.applications = regions
            .iter()
            .flat_map(|item| item.applications.clone())
            .collect::<Vec<Application>>();
        debug!(
            applications = format!("{:?}", self.applications),
            "Warden scan accomplished."
        )
    }

    pub async fn start(&self) {
        for item in self.applications.iter() {
            let mut instance = AppInstance::new();
            match instance.start(item.clone()).await {
                Ok(_) => {
                    debug!(id = item.id, "Warden successfully created instance for");
                    INSTANCES.lock().await.insert(item.clone().id, instance);
                }
                Err(err) => warn!(
                    id = item.id,
                    err = format!("{:?}", err),
                    "Warden failed to create an instance for"
                ),
            };
        }
    }
}

impl Default for WardenInstance {
    fn default() -> Self {
        Self::new()
    }
}

#[derive(Debug, Object, Clone, Serialize, Deserialize)]
pub struct Application {
    pub id: String,
    pub exe: String,
    pub args: Option<Vec<String>>,
    pub env: Option<HashMap<String, String>>,
    pub workdir: String,
}
