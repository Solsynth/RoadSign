use std::error;
use actix_web::{App, HttpServer, web};
use actix_web::dev::Server;
use actix_web::middleware::Logger;
use awc::Client;
use crate::config::CFG;
use crate::proxies::route;
use crate::server::ServerBindConfig;
use crate::tls::{load_certificates, use_rustls};

pub async fn build_proxies() -> Result<Vec<Server>, Box<dyn error::Error>> {
    load_certificates().await?;

    let cfg = CFG
        .read()
        .await
        .get::<Vec<ServerBindConfig>>("proxies.bind")?;

    let mut tasks = Vec::new();
    for item in cfg {
        tasks.push(build_single_proxy(item)?);
    }

    Ok(tasks)
}

pub fn build_single_proxy(cfg: ServerBindConfig) -> Result<Server, Box<dyn error::Error>> {
    let server = HttpServer::new(|| {
        App::new()
            .wrap(Logger::default())
            .app_data(web::Data::new(Client::default()))
            .route("/", web::to(route::handle))
    });
    if cfg.tls {
        Ok(server.bind_rustls_0_22(cfg.addr, use_rustls()?)?.run())
    } else {
        Ok(server.bind(cfg.addr)?.run())
    }
}