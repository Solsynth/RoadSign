use futures_util::{SinkExt, StreamExt};
use http::{header, request::Builder, HeaderMap, Method, StatusCode, Uri};
use lazy_static::lazy_static;
use poem::{
    web::{websocket::WebSocket, StaticFileRequest},
    Body, Error, FromRequest, IntoResponse, Request, Response,
};
use std::{
    ffi::OsStr,
    path::{Path, PathBuf},
    sync::Arc,
};
use tokio::sync::RwLock;
use tokio_tungstenite::connect_async;

use super::browser::{DirectoryTemplate, FileRef};

lazy_static! {
    pub static ref CLIENT: reqwest::Client = reqwest::Client::new();
}

pub async fn repond_websocket(req: Builder, ws: WebSocket) -> Response {
    ws.on_upgrade(move |socket| async move {
        let (mut clientsink, mut clientstream) = socket.split();

        // Start connection to server
        let (serversocket, _) = connect_async(req.body(()).unwrap()).await.unwrap();
        let (mut serversink, mut serverstream) = serversocket.split();

        let client_live = Arc::new(RwLock::new(true));
        let server_live = client_live.clone();

        tokio::spawn(async move {
            while let Some(Ok(msg)) = clientstream.next().await {
                if (serversink.send(msg.into()).await).is_err() {
                    break;
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
                if (clientsink.send(msg.into()).await).is_err() {
                    break;
                };
                if !*server_live.read().await {
                    break;
                };
            }

            *server_live.write().await = false;
        });
    })
    .into_response()
}

pub async fn respond_hypertext(
    uri: String,
    ori: &Uri,
    req: &Request,
    method: Method,
    body: Body,
    headers: &HeaderMap,
) -> Result<Response, Error> {
    let ip = req.remote_addr().to_string();
    let proto = req.uri().scheme_str().unwrap();
    let host = req.uri().host().unwrap();

    let mut headers = headers.clone();
    headers.insert("Server", "RoadSign".parse().unwrap());
    headers.insert("X-Forward-For", ip.parse().unwrap());
    headers.insert("X-Forwarded-Proto", proto.parse().unwrap());
    headers.insert("X-Forwarded-Host", host.parse().unwrap());
    headers.insert("X-Real-IP", ip.parse().unwrap());
    headers.insert(
        "Forwarded",
        format!("by={};for={};host={};proto={}", ip, ip, host, proto)
            .parse()
            .unwrap(),
    );

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
            res.headers_mut()
                .insert("Server", "RoadSign".parse().unwrap());
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

pub struct StaticResponderConfig {
    pub uri: String,
    pub utf8: bool,
    pub browse: bool,
    pub with_slash: bool,
    pub index: Option<String>,
    pub fallback: Option<String>,
    pub suffix: Option<String>,
}

pub async fn respond_static(
    cfg: StaticResponderConfig,
    method: Method,
    req: &Request,
) -> Result<Response, Error> {
    if method != Method::GET {
        return Err(Error::from_string(
            "This destination only support GET request.",
            StatusCode::METHOD_NOT_ALLOWED,
        ));
    }

    let path = req
        .uri()
        .path()
        .trim_start_matches('/')
        .trim_end_matches('/');

    let path = percent_encoding::percent_decode_str(path)
        .decode_utf8()
        .map_err(|_| Error::from_status(StatusCode::NOT_FOUND))?;

    let base_path = cfg.uri.parse::<PathBuf>().unwrap();
    let mut file_path = base_path.clone();
    for p in Path::new(&*path) {
        if p == OsStr::new(".") {
            continue;
        } else if p == OsStr::new("..") {
            file_path.pop();
        } else {
            file_path.push(p);
        }
    }

    if !file_path.starts_with(cfg.uri) {
        return Err(Error::from_status(StatusCode::FORBIDDEN));
    }

    if !file_path.exists() {
        if let Some(suffix) = cfg.suffix {
            let file_name = file_path
                .file_name()
                .and_then(OsStr::to_str)
                .unwrap()
                .to_string();
            file_path.pop();
            file_path.push((file_name + &suffix).as_str());
            if file_path.is_file() {
                return Ok(StaticFileRequest::from_request_without_body(req)
                    .await?
                    .create_response(&file_path, cfg.utf8)?
                    .into_response());
            }
        }

        if let Some(file) = cfg.fallback {
            let fallback_path = base_path.join(file);
            if fallback_path.is_file() {
                return Ok(StaticFileRequest::from_request_without_body(req)
                    .await?
                    .create_response(&fallback_path, cfg.utf8)?
                    .into_response());
            }
        }
        return Err(Error::from_status(StatusCode::NOT_FOUND));
    }

    if file_path.is_file() {
        Ok(StaticFileRequest::from_request_without_body(req)
            .await?
            .create_response(&file_path, cfg.utf8)?
            .into_response())
    } else {
        if cfg.with_slash
            && !req.original_uri().path().ends_with('/')
            && (cfg.index.is_some() || cfg.browse)
        {
            let redirect_to = format!("{}/", req.original_uri().path());
            return Ok(Response::builder()
                .status(StatusCode::FOUND)
                .header(header::LOCATION, redirect_to)
                .finish());
        }

        if let Some(index_file) = &cfg.index {
            let index_path = file_path.join(index_file);
            if index_path.is_file() {
                return Ok(StaticFileRequest::from_request_without_body(req)
                    .await?
                    .create_response(&index_path, cfg.utf8)?
                    .into_response());
            }
        }

        if cfg.browse {
            let read_dir = file_path
                .read_dir()
                .map_err(|_| Error::from_status(StatusCode::FORBIDDEN))?;
            let mut template = DirectoryTemplate {
                path: &path,
                files: Vec::new(),
            };

            for res in read_dir {
                let entry = res.map_err(|_| Error::from_status(StatusCode::FORBIDDEN))?;

                if let Some(filename) = entry.file_name().to_str() {
                    let mut base_url = req.original_uri().path().to_string();
                    if !base_url.ends_with('/') {
                        base_url.push('/');
                    }
                    let filename_url = percent_encoding::percent_encode(
                        filename.as_bytes(),
                        percent_encoding::NON_ALPHANUMERIC,
                    );
                    template.files.push(FileRef {
                        url: format!("{base_url}{filename_url}"),
                        filename: filename.to_string(),
                        is_dir: entry.path().is_dir(),
                    });
                }
            }

            let html = template.render();
            Ok(Response::builder()
                .header(header::CONTENT_TYPE, mime::TEXT_HTML_UTF_8.as_ref())
                .body(Body::from_string(html)))
        } else {
            Err(Error::from_status(StatusCode::NOT_FOUND))
        }
    }
}
