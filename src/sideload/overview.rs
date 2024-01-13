use poem_openapi::{param::Query, payload::PlainText, OpenApi};

use super::SideloadApi;

#[OpenApi]
impl SideloadApi {
    #[oai(path = "/hello", method = "get")]
    async fn index(&self, name: Query<Option<String>>) -> PlainText<String> {
        match name.0 {
            Some(name) => PlainText(format!("hello, {name}!")),
            None => PlainText("hello!".to_string()),
        }
    }
}
