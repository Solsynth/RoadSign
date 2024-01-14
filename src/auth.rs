use http::StatusCode;
use poem::{
    web::headers::{self, authorization::Basic, HeaderMapExt},
    Endpoint, Error, Middleware, Request, Response, Result,
};

pub struct BasicAuth {
    pub username: String,
    pub password: String,
}

impl<E: Endpoint> Middleware<E> for BasicAuth {
    type Output = BasicAuthEndpoint<E>;

    fn transform(&self, ep: E) -> Self::Output {
        BasicAuthEndpoint {
            ep,
            username: self.username.clone(),
            password: self.password.clone(),
        }
    }
}

pub struct BasicAuthEndpoint<E> {
    ep: E,
    username: String,
    password: String,
}

#[poem::async_trait]
impl<E: Endpoint> Endpoint for BasicAuthEndpoint<E> {
    type Output = E::Output;

    async fn call(&self, req: Request) -> Result<Self::Output> {
        if let Some(auth) = req.headers().typed_get::<headers::Authorization<Basic>>() {
            if auth.0.username() == self.username && auth.0.password() == self.password {
                return self.ep.call(req).await;
            }
        }
        Err(Error::from_response(
            Response::builder()
                .header(
                    "WWW-Authenticate",
                    "Basic realm=\"RoadSig\", charset=\"UTF-8\"",
                )
                .status(StatusCode::UNAUTHORIZED)
                .finish(),
        ))
    }
}
