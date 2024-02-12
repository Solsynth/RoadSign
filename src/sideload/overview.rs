use actix_web::web;
use serde::Serialize;
use crate::proxies::config::{Destination, Location};
use crate::proxies::metrics::RoadTrace;
use crate::ROAD;

#[derive(Debug, Clone, PartialEq, Serialize)]
pub struct OverviewData {
    regions: usize,
    locations: usize,
    destinations: usize,
    requests_count: u64,
    failures_count: u64,
    successes_count: u64,
    success_rate: f64,
    recent_successes: Vec<RoadTrace>,
    recent_errors: Vec<RoadTrace>,
}

pub async fn get_overview() -> web::Json<OverviewData> {
    let locked_app = ROAD.lock().await;
    let regions = locked_app.regions.clone();
    let locations = regions
        .iter()
        .flat_map(|item| item.locations.clone())
        .collect::<Vec<Location>>();
    let destinations = locations
        .iter()
        .flat_map(|item| item.destinations.clone())
        .collect::<Vec<Destination>>();
    web::Json(OverviewData {
        regions: regions.len(),
        locations: locations.len(),
        destinations: destinations.len(),
        requests_count: locked_app.metrics.requests_count,
        successes_count: locked_app.metrics.requests_count - locked_app.metrics.failures_count,
        failures_count: locked_app.metrics.failures_count,
        success_rate: locked_app.metrics.get_success_rate(),
        recent_successes: locked_app
            .metrics
            .recent_successes
            .clone()
            .into_iter()
            .collect::<Vec<_>>(),
        recent_errors: locked_app
            .metrics
            .recent_errors
            .clone()
            .into_iter()
            .collect::<Vec<_>>(),
    })
}