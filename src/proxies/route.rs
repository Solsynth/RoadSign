use poem::{
    handler,
    http::{HeaderMap, StatusCode, Uri},
    web::{websocket::WebSocket, Data},
    Body, Error, FromRequest, IntoResponse, Request, Response, Result,
};
use rand::seq::SliceRandom;
use reqwest::Method;

use crate::proxies::{
    config::{Destination, DestinationType},
    responder,
};

#[handler]
pub async fn handle(
    app: Data<&super::Instance>,
    req: &Request,
    uri: &Uri,
    headers: &HeaderMap,
    method: Method,
    body: Body,
) -> Result<impl IntoResponse, Error> {
    let location = match app.filter(uri, method.clone(), headers) {
        Some(val) => val,
        None => {
            return Err(Error::from_string(
                "There are no region be able to respone this request.",
                StatusCode::NOT_FOUND,
            ))
        }
    };

    let destination = location
        .destinations
        .choose_weighted(&mut rand::thread_rng(), |item| item.weight.unwrap_or(1))
        .unwrap();

    async fn forward(
        end: &Destination,
        req: &Request,
        ori: &Uri,
        headers: &HeaderMap,
        method: Method,
        body: Body,
    ) -> Result<Response, Error> {
        // Handle websocket
        if let Ok(ws) = WebSocket::from_request_without_body(req).await {
            // Get uri
            let Ok(uri) = end.get_websocket_uri() else {
                return Err(Error::from_string(
                    "This destination was not support websockets.",
                    StatusCode::NOT_IMPLEMENTED,
                ));
            };

            // Build request
            let mut ws_req = http::Request::builder().uri(&uri);
            for (key, value) in headers.iter() {
                ws_req = ws_req.header(key, value);
            }

            // Start the websocket connection
            return Ok(responder::repond_websocket(ws_req, ws).await);
        }

        // Handle normal web request
        match end.get_type() {
            DestinationType::Hypertext => {
                let Ok(uri) = end.get_hypertext_uri() else {
                    return Err(Error::from_string(
                        "This destination was not support web requests.",
                        StatusCode::NOT_IMPLEMENTED,
                    ));
                };

                responder::respond_hypertext(uri, ori, method, body, headers).await
            }
            DestinationType::StaticFiles => {
                let Ok(cfg) = end.get_static_config() else {
                    return Err(Error::from_string(
                        "This destination was not support static files.",
                        StatusCode::NOT_IMPLEMENTED,
                    ));
                };

                responder::respond_static(cfg, method, req).await
            }
            _ => Err(Error::from_string(
                "Unsupported destination protocol.",
                StatusCode::NOT_IMPLEMENTED,
            )),
        }
    }

    forward(destination, req, uri, headers, method, body).await
}
