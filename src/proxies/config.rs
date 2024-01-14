use std::collections::HashMap;

use poem_openapi::Object;
use queryst::parse;
use serde::{Deserialize, Serialize};
use serde_json::json;

use super::responder::StaticResponderConfig;

#[derive(Debug, Object, Clone, Serialize, Deserialize)]
pub struct Region {
    pub id: String,
    pub locations: Vec<Location>,
}

#[derive(Debug, Object, Clone, Serialize, Deserialize)]
pub struct Location {
    pub id: String,
    pub hosts: Vec<String>,
    pub paths: Vec<String>,
    pub headers: Option<HashMap<String, String>>,
    pub queries: Option<Vec<String>>,
    pub methods: Option<Vec<String>>,
    pub destinations: Vec<Destination>,
}

#[derive(Debug, Object, Clone, Serialize, Deserialize)]
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
        self.uri
            .as_str()
            .splitn(2, '?')
            .collect::<Vec<_>>()
            .get(1)
            .unwrap_or(&"")
    }

    pub fn get_host(&self) -> &str {
        (self
            .uri
            .as_str()
            .splitn(2, "://")
            .collect::<Vec<_>>()
            .get(1)
            .unwrap_or(&""))
        .splitn(2, '?')
        .collect::<Vec<_>>()[0]
    }

    pub fn get_websocket_uri(&self) -> Result<String, ()> {
        let parts = self.uri.as_str().splitn(2, "://").collect::<Vec<_>>();
        let url = parts.get(1).unwrap_or(&"");
        match self.get_protocol() {
            "http" | "https" => Ok(url.replace("http", "ws")),
            _ => Err(()),
        }
    }

    pub fn get_hypertext_uri(&self) -> Result<String, ()> {
        match self.get_protocol() {
            "http" => Ok("http://".to_string() + self.get_host()),
            "https" => Ok("https://".to_string() + self.get_host()),
            _ => Err(()),
        }
    }

    pub fn get_static_config(&self) -> Result<StaticResponderConfig, ()> {
        match self.get_protocol() {
            "file" | "files" => {
                let queries = parse(self.get_queries()).unwrap_or(json!({}));
                Ok(StaticResponderConfig {
                    uri: self.get_host().to_string(),
                    utf8: queries
                        .get("utf8")
                        .and_then(|val| val.as_bool())
                        .unwrap_or(false),
                    browse: queries
                        .get("browse")
                        .and_then(|val| val.as_bool())
                        .unwrap_or(false),
                    with_slash: queries
                        .get("slash")
                        .and_then(|val| val.as_bool())
                        .unwrap_or(false),
                    index: queries
                        .get("index")
                        .and_then(|val| val.as_str().map(str::to_string)),
                    fallback: queries
                        .get("fallback")
                        .and_then(|val| val.as_str().map(str::to_string)),
                    suffix: queries
                        .get("suffix")
                        .and_then(|val| val.as_str().map(str::to_string)),
                })
            }
            _ => Err(()),
        }
    }
}
