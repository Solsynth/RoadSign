extern crate core;

mod config;
mod proxies;
mod sideload;
mod warden;
mod server;
pub mod tls;

use std::error;
use lazy_static::lazy_static;
use proxies::RoadInstance;
use tokio::sync::Mutex;
use tokio::task::JoinSet;
use tracing::{error, info, Level};
use crate::proxies::server::build_proxies;
use crate::sideload::server::build_sideload;

lazy_static! {
    static ref ROAD: Mutex<RoadInstance> = Mutex::new(RoadInstance::new());
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn error::Error>> {
    // Setting up logging
    tracing_subscriber::fmt()
        .with_max_level(Level::DEBUG)
        .init();

    // Prepare all the stuff
    info!("Loading proxy regions...");
    match proxies::loader::scan_regions(
        config::CFG
            .read()
            .await
            .get_string("regions")?
    ) {
        Err(_) => error!("Loading proxy regions... failed"),
        Ok((regions, count)) => {
            ROAD.lock().await.regions = regions;
            info!(count, "Loading proxy regions... done")
        }
    };

    let mut server_set = JoinSet::new();

    // Proxies
    for server in build_proxies().await? {
        server_set.spawn(server);
    }

    // Sideload
    server_set.spawn(build_sideload().await?);

    // Process manager
    {
        let mut app = ROAD.lock().await;
        {
            let reg = app.regions.clone();
            app.warden.scan(reg);
        }
        app.warden.start().await;
    }

    // Wait for web servers
    server_set.join_next().await;

    Ok(())
}
