mod config;
mod proxies;
mod sideload;

use poem::{listener::TcpListener, Endpoint, EndpointExt, Route, Server};
use poem_openapi::OpenApiService;
use tracing::{error, info, Level};

use crate::proxies::route;

#[tokio::main]
async fn main() -> Result<(), std::io::Error> {
    // Setting up logging
    if std::env::var_os("RUST_LOG").is_none() {
        std::env::set_var("RUST_LOG", "poem=debug");
    }
    tracing_subscriber::fmt()
        .with_max_level(Level::DEBUG)
        .init();

    // Prepare all the stuff
    let mut instance = proxies::Instance::new();

    info!("Loading proxy regions...");
    match proxies::loader::scan_regions(
        config::C
            .read()
            .unwrap()
            .get_string("regions")
            .unwrap_or("./regions".to_string()),
    ) {
        Err(_) => error!("Loading proxy regions... failed"),
        Ok((regions, count)) => {
            instance.regions = regions;
            info!(count, "Loading proxy regions... done")
        }
    };

    // Proxies
    let proxies_server = Server::new(TcpListener::bind(
        config::C
            .read()
            .unwrap()
            .get_string("listen.proxies")
            .unwrap_or("0.0.0.0:80".to_string()),
    ))
    .run(route::handle.data(instance));

    // Sideload
    let sideload = OpenApiService::new(sideload::SideloadApi, "Sideload API", "1.0")
        .server("http://localhost:3000/cgi");
    let sideload_ui = sideload.swagger_ui();

    let sideload_server = Server::new(TcpListener::bind(
        config::C
            .read()
            .unwrap()
            .get_string("listen.sideload")
            .unwrap_or("0.0.0.0:81".to_string()),
    ))
    .run(
        Route::new()
            .nest("/cgi", sideload)
            .nest("/swagger", sideload_ui),
    );

    tokio::try_join!(proxies_server, sideload_server)?;

    Ok(())
}
