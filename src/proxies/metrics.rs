use std::collections::VecDeque;

use poem_openapi::Object;
use serde::{Deserialize, Serialize};

use super::config::{Destination, Location, Region};

#[derive(Debug, Object, Clone, Serialize, Deserialize, PartialEq)]
pub struct RoadTrace {
    pub region: String,
    pub location: String,
    pub destination: String,
    pub error: Option<String>,
}

impl RoadTrace {
    pub fn from_structs(reg: Region, loc: Location, end: Destination) -> RoadTrace {
        RoadTrace {
            region: reg.id,
            location: loc.id,
            destination: end.id,
            error: None,
        }
    }

    pub fn from_structs_with_error(
        reg: Region,
        loc: Location,
        end: Destination,
        err: String,
    ) -> RoadTrace {
        RoadTrace {
            region: reg.id,
            location: loc.id,
            destination: end.id,
            error: Some(err),
        }
    }
}

#[derive(Debug, Clone)]
pub struct RoadMetrics {
    pub requests_count: u64,
    pub failures_count: u64,

    pub recent_successes: VecDeque<RoadTrace>,
    pub recent_errors: VecDeque<RoadTrace>,
}

impl RoadMetrics {
    pub fn get_success_rate(&self) -> f64 {
        if self.requests_count > 0 {
            (self.requests_count - self.failures_count) as f64 / self.requests_count as f64
        } else {
            0.0
        }
    }

    pub fn add_success_request(&mut self, reg: Region, loc: Location, end: Destination) {
        self.requests_count += 1;
        self.recent_successes
            .push_back(RoadTrace::from_structs(reg, loc, end));
    }

    pub fn add_faliure_request(
        &mut self,
        reg: Region,
        loc: Location,
        end: Destination,
        err: String, // For some reason error is rarely clonable, so we use preformatted message
    ) {
        self.requests_count += 1;
        self.failures_count += 1;
        self.recent_errors
            .push_back(RoadTrace::from_structs_with_error(reg, loc, end, err));
        if self.recent_errors.len() > 10 {
            self.recent_errors.pop_front();
        }
    }
}
