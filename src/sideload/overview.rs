use poem_openapi::{payload::Json, ApiResponse, Object};

use crate::{
    proxies::{
        config::{Destination, Location},
        metrics::RoadTrace,
    },
    ROAD,
};

#[derive(ApiResponse)]
pub enum OverviewResponse {
    /// Return the overview data.
    #[oai(status = 200)]
    Ok(Json<OverviewData>),
}

#[derive(Debug, Object, Clone, PartialEq)]
pub struct OverviewData {
    /// Loaded regions count
    #[oai(read_only)]
    regions: usize,
    /// Loaded locations count
    #[oai(read_only)]
    locations: usize,
    /// Loaded destnations count
    #[oai(read_only)]
    destinations: usize,
    /// Recent requests count
    requests_count: u64,
    /// Recent requests success count
    faliures_count: u64,
    /// Recent requests falied count
    successes_count: u64,
    /// Recent requests success rate
    success_rate: f64,
    /// Recent successes
    recent_successes: Vec<RoadTrace>,
    /// Recent errors
    recent_errors: Vec<RoadTrace>,
}

pub async fn index() -> OverviewResponse {
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
    OverviewResponse::Ok(Json(OverviewData {
        regions: regions.len(),
        locations: locations.len(),
        destinations: destinations.len(),
        requests_count: locked_app.metrics.requests_count,
        successes_count: locked_app.metrics.requests_count - locked_app.metrics.failures_count,
        faliures_count: locked_app.metrics.failures_count,
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
    }))
}
