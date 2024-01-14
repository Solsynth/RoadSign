use poem_openapi::OpenApi;

pub mod overview;
pub mod regions;

pub struct SideloadApi;

#[OpenApi]
impl SideloadApi {
    #[oai(path = "/", method = "get")]
    async fn index(&self) -> overview::OverviewResponse {
        overview::index().await
    }

    #[oai(path = "/regions", method = "get")]
    async fn regions_index(&self) -> regions::RegionResponse {
        regions::index().await
    }
}
