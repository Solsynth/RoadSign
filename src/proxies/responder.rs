use futures_util::{SinkExt};
use std::{
    ffi::OsStr,
    path::{Path, PathBuf},
};
use actix_files::{NamedFile};
use actix_proxy::IntoHttpResponse;
use actix_web::{HttpRequest, HttpResponse, web};
use actix_web::http::Method;
use awc::Client;

pub async fn respond_hypertext(
    uri: String,
    req: HttpRequest,
    client: web::Data<Client>,
) -> Result<HttpResponse, HttpResponse> {
    let ip = req.peer_addr().unwrap().ip().to_string();
    let proto = req.uri().scheme_str().unwrap();
    let host = req.uri().host().unwrap();

    let mut headers = req.headers().clone();
    headers.insert("Server".parse().unwrap(), "RoadSign".parse().unwrap());
    headers.insert("X-Forward-For".parse().unwrap(), ip.parse().unwrap());
    headers.insert("X-Forwarded-Proto".parse().unwrap(), proto.parse().unwrap());
    headers.insert("X-Forwarded-Host".parse().unwrap(), host.parse().unwrap());
    headers.insert("X-Real-IP".parse().unwrap(), ip.parse().unwrap());
    headers.insert(
        "Forwarded".parse().unwrap(),
        format!("by={};for={};host={};proto={}", ip, ip, host, proto)
            .parse()
            .unwrap(),
    );

    let res = client.get(uri).send().await;

    return match res {
        Ok(result) => {
            let mut res = result.into_http_response();
            res.headers_mut().insert("Server".parse().unwrap(), "RoadSign".parse().unwrap());
            Ok(res)
        }

        Err(error) => {
            Err(HttpResponse::BadGateway()
                .body(format!("Something went wrong... {:}", error)))
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
) -> Result<HttpResponse, HttpResponse> {
    if req.method() != Method::GET {
        return Err(HttpResponse::MethodNotAllowed()
            .body("This destination only support GET request."));
    }

    let path = req
        .uri()
        .path()
        .trim_start_matches('/')
        .trim_end_matches('/');

    let path = match percent_encoding::percent_decode_str(path).decode_utf8() {
        Ok(val) => val,
        Err(_) => {
            return Err(HttpResponse::NotFound().body("Not found."));
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
        return Err(HttpResponse::Forbidden()
            .body("Unexpected path."));
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

        return Err(HttpResponse::NotFound().body("Not found."));
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

        Err(HttpResponse::NotFound().body("Not found."))
    };
}
