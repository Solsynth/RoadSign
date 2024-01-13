use config::Config;

pub fn load_settings() -> Config {
    Config::builder()
        .add_source(config::File::with_name("Settings"))
        .add_source(config::Environment::with_prefix("ROADSIGN"))
        .build()
        .unwrap()
}
