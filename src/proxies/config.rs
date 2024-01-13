use std::{collections::HashMap, time::Duration};

#[derive(Debug, Clone)]
pub struct Region {
    id: String,
    locations: Vec<Location>,
}

#[derive(Debug, Clone)]
pub struct Location {
    hosts: Vec<String>,
    paths: Vec<String>,
    headers: Vec<HashMap<String, String>>,
    query_strings: Vec<HashMap<String, String>>,
    destination: Vec<Destination>,
}

#[derive(Debug, Clone)]
pub struct Destination {
    uri: Vec<String>,
    timeout: Duration,
}
