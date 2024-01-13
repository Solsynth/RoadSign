use std::collections::HashMap;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Region {
    id: String,
    locations: Vec<Location>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Location {
    id: String,
    hosts: Vec<String>,
    paths: Vec<String>,
    headers: Option<Vec<HashMap<String, String>>>,
    query_strings: Option<Vec<HashMap<String, String>>>,
    destinations: Vec<Destination>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Destination {
    id: String,
    uri: String,
    timeout: Option<u32>,
}
