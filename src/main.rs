mod config;
mod proxies;
mod sideload;
pub mod warden;
mod tls;

use std::error;
use actix_web::{App, HttpServer, web};
use actix_web::middleware::Logger;
use actix_web_httpauth::extractors::AuthenticationError;
use actix_web_httpauth::headers::www_authenticate::basic::Basic;
use actix_web_httpauth::middleware::HttpAuthentication;
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

    // Proxies
    let proxies_server = HttpServer::new(|| {
        App::new()
            .wrap(Logger::default())
            .app_data(web::Data::new(Client::default()))
            .route("/", web::to(route::handle))
    }).bind_rustls_0_22(
        config::CFG
            .read()
            .await
            .get_string("listen.proxies_tls")?,
        tls::use_rustls().await?,
    )?.bind(
        config::CFG
            .read()
            .await
            .get_string("listen.proxies")?
    )?.run();

    // Sideload
    let sideload_server = HttpServer::new(|| {
        App::new()
            .wrap(HttpAuthentication::basic(|req, credentials| async move {
                let password = match config::CFG
                    .read()
                    .await
                    .get_string("secret") {
                    Ok(val) => val,
                    Err(_) => return Err((AuthenticationError::new(Basic::new()).into(), req))
                };
                if credentials.password().unwrap_or("") != password {
                    Err((AuthenticationError::new(Basic::new()).into(), req))
                } else {
                    Ok(req)
                }
            }))
            .service(sideload::service())
    }).bind(
        config::CFG
            .read()
            .await
            .get_string("listen.sideload")
            .unwrap_or("0.0.0.0:81".to_string())
    )?.workers(1).run();

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
