use std::{
    ffi::OsStr,
    path::{Path, PathBuf},
};
use actix_files::{NamedFile};
use actix_proxy::IntoHttpResponse;
use actix_web::{HttpRequest, HttpResponse, web};
use actix_web::http::{header, Method};
use actix_web::http::header::HeaderValue;
use awc::Client;
use tracing::log::warn;
use crate::proxies::ProxyError;

pub async fn respond_hypertext(
    uri: String,
    req: HttpRequest,
    client: web::Data<Client>,
) -> Result<HttpResponse, ProxyError> {
    let conn = req.connection_info();
    let ip = conn.realip_remote_addr().unwrap_or("0.0.0.0");
    let proto = conn.scheme();
    let host = conn.host();

    let mut headers = req.headers().clone();
    headers.insert(header::X_FORWARDED_FOR, ip.parse().unwrap());
    headers.insert(header::X_FORWARDED_PROTO, proto.parse().unwrap());
    headers.insert(header::X_FORWARDED_HOST, host.parse().unwrap());
    headers.insert(
        header::FORWARDED,
        format!("by={};for={};host={};proto={}", ip, ip, host, proto)
            .parse()
            .unwrap(),
    );

    let append_part = req.uri().to_string().chars().skip(1).collect::<String>();
    let target_url = format!("{}{}", uri, append_part);

    let res = client.request(req.method().clone(), target_url).send().await;

    return match res {
        Ok(result) => {
            let mut res = result.into_http_response();
            res.headers_mut().insert(header::SERVER, HeaderValue::from_static("RoadSign"));
            res.headers_mut().remove(header::CONTENT_ENCODING);
            Ok(res)
        }

        Err(error) => {
            warn!("Proxy got a upstream issue... {:?}", error);
            Err(ProxyError::BadGateway)
        }
    };
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

    return if file_path.is_file() {
        Ok(NamedFile::open(file_path).unwrap().into_response(&req))
    } else {
        if let Some(index_file) = &cfg.index {
            let index_path = file_path.join(index_file);
            if index_path.is_file() {
                return Ok(NamedFile::open(index_path).unwrap().into_response(&req));
            }
        }

        return Err(ProxyError::NotFound);
    };
}
