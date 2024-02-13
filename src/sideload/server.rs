use std::error;
use actix_web::dev::Server;
use actix_web::{App, HttpServer};
use actix_web_httpauth::extractors::AuthenticationError;
use actix_web_httpauth::headers::www_authenticate::basic::Basic;
use actix_web_httpauth::middleware::HttpAuthentication;
use crate::sideload;

pub async fn build_sideload() -> Result<Server, Box<dyn error::Error>> {
    Ok(
        HttpServer::new(|| {
            App::new()
                .wrap(HttpAuthentication::basic(|req, credentials| async move {
                    let password = match crate::config::CFG
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
            crate::config::CFG
                .read()
                .await
                .get_string("sideload.bind_addr")?
        )?.workers(1).run()
    )
}