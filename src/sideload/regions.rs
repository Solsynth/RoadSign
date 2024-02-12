use actix_web::web;
use crate::proxies::config::Region;
use crate::ROAD;

pub async fn list_region() -> web::Json<Vec<Region>> {
    let locked_app = ROAD.lock().await;

    web::Json(locked_app.regions.clone())
}