use futures_util::{SinkExt, StreamExt};
use poem::{
    handler,
    http::{HeaderMap, StatusCode, Uri},
    web::{websocket::WebSocket, Data},
    Body, Error, FromRequest, IntoResponse, Request, Response, Result,
};
use rand::seq::SliceRandom;
use reqwest::Method;

use lazy_static::lazy_static;
use std::sync::Arc;
use tokio::sync::RwLock;
use tokio_tungstenite::connect_async;

use crate::proxies::config::{Destination, DestinationType};

lazy_static! {
    pub static ref CLIENT: reqwest::Client = reqwest::Client::new();
}

#[handler]
pub async fn handle(
    app: Data<&super::Instance>,
    req: &Request,
    uri: &Uri,
    headers: &HeaderMap,
    method: Method,
    body: Body,
) -> Result<impl IntoResponse, Error> {
    let location = match app.filter(uri, headers) {
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
    ) -> Result<impl IntoResponse, Error> {
        // Handle websocket
        if let Ok(ws) = WebSocket::from_request_without_body(req).await {
            // Get uri
            let Ok(uri) = end.get_websocket_uri() else {
                return Err(Error::from_string(
                    "Proxy endpoint not configured to support websockets!",
                    StatusCode::NOT_IMPLEMENTED,
                ));
            };

            // Build request
            let mut ws_req = http::Request::builder().uri(&uri);
            for (key, value) in headers.iter() {
                ws_req = ws_req.header(key, value);
            }

            // Start the websocket connection
            return Ok(ws
                .on_upgrade(move |socket| async move {
                    let (mut clientsink, mut clientstream) = socket.split();

                    // Start connection to server
                    let (serversocket, _) = connect_async(ws_req.body(()).unwrap()).await.unwrap();
                    let (mut serversink, mut serverstream) = serversocket.split();

                    let client_live = Arc::new(RwLock::new(true));
                    let server_live = client_live.clone();

                    tokio::spawn(async move {
                        while let Some(Ok(msg)) = clientstream.next().await {
                            match serversink.send(msg.into()).await {
                                Err(_) => break,
                                _ => {}
                            };
                            if !*client_live.read().await {
                                break;
                            };
                        }

                        *client_live.write().await = false;
                    });

                    // Relay server messages to the client
                    tokio::spawn(async move {
                        while let Some(Ok(msg)) = serverstream.next().await {
                            match clientsink.send(msg.into()).await {
                                Err(_) => break,
                                _ => {}
                            };
                            if !*server_live.read().await {
                                break;
                            };
                        }

                        *server_live.write().await = false;
                    });
                })
                .into_response());
        }

        // Handle normal web request
        match end.get_type() {
            DestinationType::Hypertext => {
                let Ok(uri) = end.get_hypertext_uri() else {
                    return Err(Error::from_string(
                        "Proxy endpoint not configured to support web requests!",
                        StatusCode::NOT_IMPLEMENTED,
                    ));
                };

                let res = CLIENT
                    .request(method, uri + ori.path() + ori.query().unwrap_or(""))
                    .headers(headers.clone())
                    .body(body.into_bytes().await.unwrap())
                    .send()
                    .await;

                match res {
                    Ok(result) => {
                        let mut res = Response::default();
                        res.extensions().clone_from(&result.extensions());
                        result.headers().iter().for_each(|(key, val)| {
                            res.headers_mut().insert(key, val.to_owned());
                        });
                        res.set_status(result.status());
                        res.set_version(result.version());
                        res.set_body(result.bytes().await.unwrap());
                        Ok(res)
                    }

                    Err(error) => Err(Error::from_string(
                        error.to_string(),
                        error.status().unwrap_or(StatusCode::BAD_GATEWAY),
                    )),
                }
            }
            _ => Err(Error::from_status(StatusCode::NOT_IMPLEMENTED)),
        }
    }

    forward(destination, req, uri, headers, method, body).await
}
