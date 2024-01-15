pub mod auth;
mod config;
mod proxies;
mod sideload;
pub mod warden;

use lazy_static::lazy_static;
use poem::{listener::TcpListener, EndpointExt, Route, Server};
use poem_openapi::OpenApiService;
use proxies::RoadInstance;
use tokio::sync::Mutex;
use tracing::{error, info, Level};

use crate::proxies::route;

lazy_static! {
    static ref ROAD: Mutex<RoadInstance> = Mutex::new(RoadInstance::new());
}

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
    info!("Loading proxy regions...");
    match proxies::loader::scan_regions(
        config::C
            .read()
            .await
            .get_string("regions")
            .unwrap_or("./regions".to_string()),
    ) {
        Err(_) => error!("Loading proxy regions... failed"),
        Ok((regions, count)) => {
            ROAD.lock().await.regions = regions;
            info!(count, "Loading proxy regions... done")
        }
    };

    // Proxies
    let proxies_server = Server::new(TcpListener::bind(
        config::C
            .read()
            .await
            .get_string("listen.proxies")
            .unwrap_or("0.0.0.0:80".to_string()),
    ))
    .run(route::handle);

    // Sideload
    let sideload = OpenApiService::new(sideload::SideloadApi, "Sideload API", "1.0")
        .server("http://localhost:3000/cgi");

    let sideload_server = Server::new(TcpListener::bind(
        config::C
            .read()
            .await
            .get_string("listen.sideload")
            .unwrap_or("0.0.0.0:81".to_string()),
    ))
    .run(
        Route::new().nest("/cgi", sideload).with(auth::BasicAuth {
            username: "RoadSign".to_string(),
            password: config::C
                .read()
                .await
                .get_string("secret")
                .unwrap_or("password".to_string()),
        }),
    );

    // Process manager
    {
        let mut app = ROAD.lock().await;
        {
            let reg = app.regions.clone();
            app.warden.scan(reg);
        }
        app.warden.start().await;
    }

    tokio::try_join!(proxies_server, sideload_server)?;

    Ok(())
}
