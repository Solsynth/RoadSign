mod config;
mod proxies;
mod sideload;

use poem::{listener::TcpListener, Route, Server};
use poem_openapi::OpenApiService;

#[tokio::main]
async fn main() -> Result<(), std::io::Error> {
    // Load settings
    let settings = config::loader::load_settings();

    println!(
        "Will listen at {:?}",
        settings.get_array("listen.proxies").unwrap()
    );

    if std::env::var_os("RUST_LOG").is_none() {
        std::env::set_var("RUST_LOG", "poem=debug");
    }
    tracing_subscriber::fmt::init();

    // Proxies

    // Sideload
    let sideload = OpenApiService::new(sideload::SideloadApi, "Sideload API", "1.0")
        .server("http://localhost:3000/cgi");

    let sideload_server = Server::new(TcpListener::bind(
        settings
            .get_string("listen.sideload")
            .unwrap_or("0.0.0.0:81".to_string()),
    ))
    .run(Route::new().nest("/cgi", sideload));

    tokio::try_join!(sideload_server)?;

    Ok(())
}
