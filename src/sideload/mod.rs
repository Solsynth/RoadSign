use actix_web::{Scope, web};
use crate::sideload::overview::get_overview;
use crate::sideload::regions::list_region;

mod overview;
mod regions;
pub mod server;

static ROOT: &str = "";

pub fn service() -> Scope {
    web::scope("/cgi")
        .route(ROOT, web::get().to(get_overview))
        .route("/regions", web::get().to(list_region))
}