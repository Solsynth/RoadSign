use std::sync::RwLock;

use config::Config;
use lazy_static::lazy_static;

use crate::config::loader::load_settings;

pub mod loader;

lazy_static! {
    pub static ref C: RwLock<Config> = RwLock::new(load_settings());
}
