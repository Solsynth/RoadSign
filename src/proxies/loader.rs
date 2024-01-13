use std::ffi::OsStr;
use std::fs::{self, DirEntry};
use std::io;

use tracing::warn;

use crate::proxies::config;

pub fn scan_regions(basepath: String) -> io::Result<(Vec<config::Region>, u32)> {
    let mut count: u32 = 0;
    let mut result = vec![];
    for entry in fs::read_dir(basepath)? {
        if let Ok(val) = load_region(entry.unwrap()) {
            result.push(val);
            count += 1;
        };
    }

    Ok((result, count))
}

pub fn load_region(file: DirEntry) -> Result<config::Region, String> {
    if file.path().extension().and_then(OsStr::to_str).unwrap() != "toml" {
        return Err("File entry wasn't toml file".to_string());
    }

    let fp = file.path();
    let content = match fs::read_to_string(fp.clone()) {
        Ok(val) => val,
        Err(err) => {
            warn!(
                err = format!("{:?}", err),
                filepath = fp.clone().to_str(),
                "An error occurred when loading region, skipped."
            );
            return Err("Failed to load file".to_string());
        }
    };

    let data: config::Region = match toml::from_str(&content) {
        Ok(val) => val,
        Err(err) => {
            warn!(
                err = format!("{:?}", err),
                filepath = fp.clone().to_str(),
                "An error occurred when parsing region, skipped."
            );
            return Err("Failed to parse file".to_string());
        }
    };

    Ok(data)
}
