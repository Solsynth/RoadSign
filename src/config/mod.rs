use config::Config;
use lazy_static::lazy_static;
use tokio::sync::RwLock;

use crate::config::loader::load_settings;

pub mod loader;

lazy_static! {
    pub static ref CFG: RwLock<Config> = RwLock::new(load_settings());
}