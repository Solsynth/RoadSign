use std::fmt::Write;

pub struct DirectoryTemplate<'a> {
    pub path: &'a str,
    pub files: Vec<FileRef>,
}

impl<'a> DirectoryTemplate<'a> {
    pub fn render(&self) -> String {
        let mut s = format!(
            r#"
        <html>
            <head>
            <title>Index of {}</title>
        </head>
        <body>
        <h1>Index of /{}</h1>
        <ul>"#,
            self.path, self.path
        );

        for file in &self.files {
            if file.is_dir {
                let _ = write!(
                    s,
                    r#"<li><a href="{}">{}/</a></li>"#,
                    file.url, file.filename
                );
            } else {
                let _ = write!(
                    s,
                    r#"<li><a href="{}">{}</a></li>"#,
                    file.url, file.filename
                );
            }
        }

        s.push_str(
            r#"</ul>
        </body>
        </html>"#,
        );

        s
    }
}

pub struct FileRef {
    pub url: String,
    pub filename: String,
    pub is_dir: bool,
}
