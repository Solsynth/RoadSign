use poem::http::{HeaderMap, Uri};
use regex::Regex;
use serde::{Deserialize, Serialize};
use wildmatch::WildMatch;

use self::config::{Location, Region};

pub mod config;
pub mod loader;
pub mod route;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Instance {
    pub regions: Vec<Region>,
}

impl Instance {
    pub fn new() -> Instance {
        Instance { regions: vec![] }
    }

    pub fn filter(&self, uri: &Uri, headers: &HeaderMap) -> Option<&Location> {
        self.regions.iter().find_map(|region| {
            region.locations.iter().find(|location| {
                let mut hosts = location.hosts.iter();
                if !hosts.any(|item| {
                    WildMatch::new(item.as_str()).matches(uri.host().unwrap_or("localhost"))
                }) {
                    return false;
                }

                let mut paths = location.paths.iter();
                if !paths.any(|item| {
                    uri.path().starts_with(item)
                        || Regex::new(item.as_str()).unwrap().is_match(uri.path())
                }) {
                    return false;
                }

                if let Some(val) = location.headers.clone() {
                    match !val.keys().all(|item| {
                        headers.get(item).unwrap()
                            == location.headers.clone().unwrap().get(item).unwrap()
                    }) {
                        true => return false,
                        false => (),
                    }
                };

                if let Some(val) = location.queries.clone() {
                    let queries: Vec<&str> = uri.query().unwrap_or("").split('&').collect();
                    if !val.iter().all(|item| queries.contains(&item.as_str())) {
                        return false;
                    }
                }

                true
            })
        })
    }
}
