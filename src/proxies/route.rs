use actix_web::{HttpRequest, HttpResponse, web};
use actix_web::http::header;
use awc::Client;
use rand::seq::SliceRandom;

use crate::{
    proxies::{
        config::{Destination, DestinationType},
        responder,
    },
    ROAD,
};

pub async fn handle(req: HttpRequest, client: web::Data<Client>) -> HttpResponse {
    let readable_app = ROAD.lock().await;
    let (region, location) = match readable_app.filter(req.uri(), req.method(), req.headers()) {
        Some(val) => val,
        None => {
            return HttpResponse::NotFound()
                .body("There are no region be able to respone this request.");
        }
    };

    let destination = location
        .destinations
        .choose_weighted(&mut rand::thread_rng(), |item| item.weight.unwrap_or(1))
        .unwrap();

    async fn forward(
        end: &Destination,
        req: HttpRequest,
        client: web::Data<Client>,
    ) -> Result<HttpResponse, HttpResponse> {
        // Handle normal web request
        match end.get_type() {
            DestinationType::Hypertext => {
                let Ok(uri) = end.get_hypertext_uri() else {
                    return Err(HttpResponse::NotImplemented()
                        .body("This destination was not support web requests."));
                };

                responder::respond_hypertext(uri, req, client).await
            }
            DestinationType::StaticFiles => {
                let Ok(cfg) = end.get_static_config() else {
                    return Err(HttpResponse::NotImplemented()
                        .body("This destination was not support static files."));
                };

                responder::respond_static(cfg, req).await
            }
            _ => {
                return Err(HttpResponse::NotImplemented()
                    .body("Unsupported destination protocol."));
            }
        }
    }

    let reg = region.clone();
    let loc = location.clone();
    let end = destination.clone();

    let ip = match req.peer_addr() {
        None => "unknown".to_string(),
        Some(val) => val.ip().to_string()
    };
    let ua = match req.headers().get(header::USER_AGENT) {
        None => "unknown".to_string(),
        Some(val) => val.to_str().unwrap().to_string(),
    };

    match forward(&end, req, client).await {
        Ok(resp) => {
            tokio::spawn(async move {
                let writable_app = &mut ROAD.lock().await;
                writable_app.metrics.add_success_request(ip, ua, reg, loc, end);
            });
            resp
        }
        Err(resp) => {
            tokio::spawn(async move {
                let writable_app = &mut ROAD.lock().await;
                writable_app
                    .metrics
                    .add_failure_request(ip, ua, reg, loc, end, "TODO".to_owned());
            });
            resp
        }
    }
}
