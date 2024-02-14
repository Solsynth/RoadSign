use crate::proxies::ProxyError;
use crate::proxies::ProxyError::{BadGateway, UpgradeRequired};
use actix_files::NamedFile;
use actix_web::http::{header, Method};
use actix_web::web::BytesMut;
use actix_web::{web, Error, HttpRequest, HttpResponse};
use awc::error::HeaderValue;
use awc::http::Uri;
use awc::Client;
use futures::{channel::mpsc::unbounded, Sink, sink::SinkExt, stream::StreamExt};
use std::str::FromStr;
use std::time::Duration;
use std::{
    ffi::OsStr,
    path::{Path, PathBuf},
};
use actix::io::{SinkWrite, WriteHandler};
use actix::{Actor, ActorContext, AsyncContext, StreamHandler};
use actix_web_actors::ws;
use actix_web_actors::ws::{CloseReason, handshake, ProtocolError, WebsocketContext};
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tracing::log::warn;

pub async fn respond_hypertext(
    uri: String,
    req: HttpRequest,
    payload: web::Payload,
    client: web::Data<Client>,
) -> Result<HttpResponse, ProxyError> {
    let mut append_part = req.uri().to_string();
    if let Some(stripped_uri) = append_part.strip_prefix('/') {
        append_part = stripped_uri.to_string();
    }

    let uri = Uri::from_str(uri.as_str()).expect("Invalid upstream");
    let target_url = format!("{}{}", uri, append_part);

    let forwarded_req = client
        .request_from(target_url.as_str(), req.head())
        .insert_header((header::HOST, uri.host().expect("Invalid upstream")));

    let forwarded_req = match req.connection_info().realip_remote_addr() {
        Some(addr) => forwarded_req
            .insert_header((header::X_FORWARDED_FOR, addr))
            .insert_header((header::X_FORWARDED_PROTO, req.connection_info().scheme()))
            .insert_header((header::X_FORWARDED_HOST, req.connection_info().host()))
            .insert_header((
                header::FORWARDED,
                format!(
                    "by={};for={};host={};proto={}",
                    addr,
                    addr,
                    req.connection_info().host(),
                    req.connection_info().scheme()
                ),
            )),
        None => forwarded_req,
    };

    if req
        .headers()
        .get(header::UPGRADE)
        .unwrap_or(&HeaderValue::from_static(""))
        .to_str()
        .unwrap_or("")
        .to_lowercase()
        == "websocket"
    {
        let uri = uri.to_string().replacen("http", "ws", 1);
        return respond_websocket(uri, req, payload).await;
    }

    let res = forwarded_req
        .timeout(Duration::from_secs(1800))
        .send_stream(payload)
        .await
        .map_err(|err| {
            warn!("Remote gateway issue... {}", err);
            BadGateway
        })?;

    let mut client_resp = HttpResponse::build(res.status());
    for (header_name, header_value) in res
        .headers()
        .iter()
        .filter(|(h, _)| *h != header::CONNECTION && *h != header::CONTENT_ENCODING)
    {
        client_resp.insert_header((header_name.clone(), header_value.clone()));
    }

    Ok(client_resp.streaming(res))
}

pub struct WebsocketProxy<S>
    where
        S: Unpin + Sink<ws::Message>,
{
    send: SinkWrite<ws::Message, S>,
}

impl<S> WriteHandler<ProtocolError> for WebsocketProxy<S>
    where
        S: Unpin + 'static + Sink<ws::Message>,
{
    fn error(&mut self, err: ProtocolError, ctx: &mut Self::Context) -> actix::Running {
        self.error(err, ctx);
        actix::Running::Stop
    }
}

impl<S> Actor for WebsocketProxy<S>
    where
        S: Unpin + 'static + Sink<ws::Message>,
{
    type Context = WebsocketContext<Self>;
}

impl<S> StreamHandler<Result<ws::Frame, ProtocolError>> for WebsocketProxy<S>
    where
        S: Unpin + Sink<ws::Message> + 'static,
{
    fn handle(&mut self, item: Result<ws::Frame, ProtocolError>, ctx: &mut Self::Context) {
        let frame = match item {
            Ok(frame) => frame,
            Err(err) => return self.error(err, ctx),
        };
        let msg = match frame {
            ws::Frame::Text(t) => match t.try_into() {
                Ok(t) => ws::Message::Text(t),
                Err(e) => {
                    self.error(e, ctx);
                    return;
                }
            },
            ws::Frame::Binary(b) => ws::Message::Binary(b),
            ws::Frame::Continuation(c) => ws::Message::Continuation(c),
            ws::Frame::Ping(p) => ws::Message::Ping(p),
            ws::Frame::Pong(p) => ws::Message::Pong(p),
            ws::Frame::Close(r) => ws::Message::Close(r),
        };

        ctx.write_raw(msg)
    }
}

impl<S> StreamHandler<Result<ws::Message, ProtocolError>> for WebsocketProxy<S>
    where
        S: Unpin + Sink<ws::Message> + 'static,
{
    fn handle(&mut self, item: Result<ws::Message, ProtocolError>, ctx: &mut Self::Context) {
        let msg = match item {
            Ok(msg) => msg,
            Err(err) => return self.error(err, ctx),
        };

        let _ = self.send.write(msg);
    }
}

impl<S> WebsocketProxy<S>
    where
        S: Unpin + Sink<ws::Message> + 'static,
{
    fn error<E>(&mut self, err: E, ctx: &mut <Self as Actor>::Context)
        where
            E: std::error::Error,
    {
        let reason = Some(CloseReason {
            code: ws::CloseCode::Error,
            description: Some(err.to_string()),
        });

        ctx.close(reason.clone());
        let _ = self.send.write(ws::Message::Close(reason));
        self.send.close();

        ctx.stop();
    }
}

pub async fn respond_websocket(
    uri: String,
    req: HttpRequest,
    payload: web::Payload,
) -> Result<HttpResponse, ProxyError> {
    let mut res = handshake(&req).map_err(|_| UpgradeRequired)?;

    let (_, conn) = awc::Client::new()
        .ws(uri)
        .connect()
        .await
        .map_err(|_| BadGateway)?;

    let (send, recv) = conn.split();

    let out = WebsocketContext::with_factory(payload, |ctx| {
        ctx.add_stream(recv);
        WebsocketProxy {
            send: SinkWrite::new(send, ctx),
        }
    });

    Ok(res.streaming(out))
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
    req: HttpRequest,
) -> Result<HttpResponse, ProxyError> {
    if req.method() != Method::GET {
        return Err(ProxyError::MethodGetOnly);
    }

    let path = req
        .uri()
        .path()
        .trim_start_matches('/')
        .trim_end_matches('/');

    let path = match percent_encoding::percent_decode_str(path).decode_utf8() {
        Ok(val) => val,
        Err(_) => {
            return Err(ProxyError::NotFound);
        }
    };

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
        return Err(ProxyError::InvalidRequestPath);
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
                return Ok(NamedFile::open(file_path).unwrap().into_response(&req));
            }
        }

        if let Some(file) = cfg.fallback {
            let fallback_path = base_path.join(file);
            if fallback_path.is_file() {
                return Ok(NamedFile::open(fallback_path).unwrap().into_response(&req));
            }
        }

        return Err(ProxyError::NotFound);
    }

    if file_path.is_file() {
        Ok(NamedFile::open(file_path).unwrap().into_response(&req))
    } else {
        if let Some(index_file) = &cfg.index {
            let index_path = file_path.join(index_file);
            if index_path.is_file() {
                return Ok(NamedFile::open(index_path).unwrap().into_response(&req));
            }
        }

        Err(ProxyError::NotFound)
    }
}
