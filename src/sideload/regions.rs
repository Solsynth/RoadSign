use poem_openapi::{payload::Json, ApiResponse};

use crate::{proxies::config::Region, ROAD};

#[derive(ApiResponse)]
pub enum RegionResponse {
    /// Return the region data.
    #[oai(status = 200)]
    Ok(Json<Region>),
    /// Return the list of region data.
    #[oai(status = 200)]
    OkMany(Json<Vec<Region>>),
    /// Return the region data after created.
    #[oai(status = 201)]
    Created(Json<Region>),
    /// Return was not found.
    #[oai(status = 404)]
    NotFound,
}

pub async fn index() -> RegionResponse {
    let locked_app = ROAD.lock().await;

    RegionResponse::OkMany(Json(locked_app.regions.clone()))
}
