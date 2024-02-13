use std::collections::VecDeque;

use serde::{Deserialize, Serialize};

use super::config::{Destination, Location, Region};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub struct RoadTrace {
    pub region: String,
    pub location: String,
    pub destination: String,
    pub ip_address: String,
    pub user_agent: String,
    pub error: Option<String>,
}

impl RoadTrace {
    pub fn from_structs(
        ip: String,
        ua: String,
        reg: Region,
        loc: Location,
        end: Destination,
    ) -> RoadTrace {
        RoadTrace {
            ip_address: ip,
            user_agent: ua,
            region: reg.id,
            location: loc.id,
            destination: end.id,
            error: None,
        }
    }

    pub fn from_structs_with_error(
        ip: String,
        ua: String,
        reg: Region,
        loc: Location,
        end: Destination,
        err: String,
    ) -> RoadTrace {
        let mut trace = Self::from_structs(ip, ua, reg, loc, end);
        trace.error = Some(err);
        return trace;
    }
}

#[derive(Debug, Clone)]
pub struct RoadMetrics {
    pub requests_count: u64,
    pub failures_count: u64,

    pub recent_successes: VecDeque<RoadTrace>,
    pub recent_errors: VecDeque<RoadTrace>,
}

const MAX_TRACE_COUNT: usize = 32;

impl RoadMetrics {
    pub fn new() -> RoadMetrics {
        RoadMetrics {
            requests_count: 0,
            failures_count: 0,
            recent_successes: VecDeque::new(),
            recent_errors: VecDeque::new(),
        }
    }

    pub fn get_success_rate(&self) -> f64 {
        if self.requests_count > 0 {
            (self.requests_count - self.failures_count) as f64 / self.requests_count as f64
        } else {
            0.0
        }
    }

    pub fn add_success_request(
        &mut self,
        ip: String,
        ua: String,
        reg: Region,
        loc: Location,
        end: Destination,
    ) {
        self.requests_count += 1;
        self.recent_successes
            .push_back(RoadTrace::from_structs(ip, ua, reg, loc, end));
        if self.recent_successes.len() > MAX_TRACE_COUNT {
            self.recent_successes.pop_front();
        }
    }

    pub fn add_failure_request(
        &mut self,
        ip: String,
        ua: String,
        reg: Region,
        loc: Location,
        end: Destination,
        err: String, // For some reason error is rarely cloneable, so we use preformatted message
    ) {
        self.requests_count += 1;
        self.failures_count += 1;
        self.recent_errors
            .push_back(RoadTrace::from_structs_with_error(ip, ua, reg, loc, end, err));
        if self.recent_errors.len() > MAX_TRACE_COUNT {
            self.recent_errors.pop_front();
        }
    }
}
