use std::fs::File;
use std::{error};
use std::io::BufReader;
use std::sync::Arc;
use config::ConfigError;
use lazy_static::lazy_static;
use rustls::crypto::ring::sign::RsaSigningKey;
use rustls::server::{ClientHello, ResolvesServerCert};
use rustls::sign::CertifiedKey;
use serde::{Deserialize, Serialize};
use std::sync::Mutex;
use wildmatch::WildMatch;

lazy_static! {
    static ref CERTS: Mutex<Vec<CertificateConfig>> = Mutex::new(Vec::new());
}

#[derive(Debug)]
struct ProxyCertResolver;

impl ResolvesServerCert for ProxyCertResolver {
    fn resolve(&self, handshake: ClientHello) -> Option<Arc<CertifiedKey>> {
        let domain = handshake.server_name()?;

        let certs = CERTS.lock().unwrap();
        for cert in certs.iter() {
            if WildMatch::new(cert.domain.as_str()).matches(domain) {
                return match cert.clone().load() {
                    Ok(val) => Some(val),
                    Err(_) => None
                };
            }
        }
        None
    }
}

#[derive(Clone, Serialize, Deserialize)]
struct CertificateConfig {
    pub domain: String,
    pub certs: String,
    pub key: String,
}

impl CertificateConfig {
    pub fn load(self) -> Result<Arc<CertifiedKey>, Box<dyn error::Error>> {
        let certs =
            rustls_pemfile::certs(&mut BufReader::new(&mut File::open(self.certs)?))
                .collect::<Result<Vec<_>, _>>()?;
        let key =
            rustls_pemfile::private_key(&mut BufReader::new(&mut File::open(self.key)?))?
                .unwrap();
        let sign = RsaSigningKey::new(&key)?;

        Ok(Arc::new(CertifiedKey::new(certs, Arc::new(sign))))
    }
}

pub async fn load_certificates() -> Result<(), ConfigError> {
    let certs = crate::config::CFG
        .read()
        .await
        .get::<Vec<CertificateConfig>>("certificates")?;

    CERTS.lock().unwrap().clone_from(&certs);

    Ok(())
}

pub fn use_rustls() -> Result<rustls::ServerConfig, ConfigError> {
    Ok(
        rustls::ServerConfig::builder()
            .with_no_client_auth()
            .with_cert_resolver(Arc::new(ProxyCertResolver))
    )
}