use std::collections::HashMap;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Region {
    pub id: String,
    pub locations: Vec<Location>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Location {
    pub id: String,
    pub hosts: Vec<String>,
    pub paths: Vec<String>,
    pub headers: Option<HashMap<String, String>>,
    pub queries: Option<Vec<String>>,
    pub destinations: Vec<Destination>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Destination {
    pub id: String,
    pub uri: String,
    pub timeout: Option<u32>,
    pub weight: Option<u32>,
}

pub enum DestinationType {
    Hypertext,
    StaticFiles,
    Unknown,
}

impl Destination {
    pub fn get_type(&self) -> DestinationType {
        match self.get_protocol() {
            "http" | "https" => DestinationType::Hypertext,
            "file" | "files" => DestinationType::StaticFiles,
            _ => DestinationType::Unknown,
        }
    }

    pub fn get_protocol(&self) -> &str {
        self.uri.as_str().splitn(2, "://").collect::<Vec<_>>()[0]
    }

    pub fn get_queries(&self) -> &str {
        self.uri.as_str().splitn(2, "?").collect::<Vec<_>>()[1]
    }

    pub fn get_host(&self) -> &str {
        (self.uri.as_str().splitn(2, "://").collect::<Vec<_>>()[1])
            .splitn(2, "?")
            .collect::<Vec<_>>()[0]
    }

    pub fn get_hypertext_uri(&self) -> Result<String, ()> {
        match self.get_protocol() {
            "http" => Ok("http://".to_string() + self.get_host()),
            "https" => Ok("https://".to_string() + self.get_host()),
            _ => Err(()),
        }
    }

    pub fn get_websocket_uri(&self) -> Result<String, ()> {
        let url = self.uri.as_str().splitn(2, "://").collect::<Vec<_>>()[1];
        match self.get_protocol() {
            "http" | "https" => Ok(url.replace("http", "ws")),
            _ => Err(()),
        }
    }
}
