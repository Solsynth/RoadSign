pub mod auth;
mod config;
mod proxies;
pub mod warden;

use actix_web::{App, HttpServer, web};
use awc::Client;
use lazy_static::lazy_static;
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
    let proxies_server = HttpServer::new(|| {
        App::new()
            .app_data(web::Data::new(Client::default()))
            .route("/", web::to(route::handle))
    }).bind(
        config::C
            .read()
            .await
            .get_string("listen.proxies")
            .unwrap_or("0.0.0.0:80".to_string())
    )?.run();

    // Process manager
    {
        let mut app = ROAD.lock().await;
        {
            let reg = app.regions.clone();
            app.warden.scan(reg);
        }
        app.warden.start().await;
    }

    tokio::try_join!(proxies_server)?;

    Ok(())
}
