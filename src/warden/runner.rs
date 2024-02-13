use std::{borrow::BorrowMut, collections::HashMap, io};

use super::Application;
use lazy_static::lazy_static;
use tokio::{
    io::{AsyncBufReadExt, BufReader},
    process::{Child, Command},
};
use tokio::sync::Mutex;

lazy_static! {
    static ref STDOUT: Mutex<HashMap<String, String>> = Mutex::new(HashMap::new());
    static ref STDERR: Mutex<HashMap<String, String>> = Mutex::new(HashMap::new());
}

#[derive(Debug)]
pub struct AppInstance {
    pub app: Option<Application>,
    pub program: Option<Child>,
}

impl AppInstance {
    pub fn new() -> AppInstance {
        AppInstance {
            app: None,
            program: None,
        }
    }

    pub async fn start(&mut self, app: Application) -> io::Result<()> {
        return match Command::new(app.exe.clone())
            .args(app.args.clone().unwrap_or_default())
            .envs(app.env.clone().unwrap_or_default())
            .current_dir(app.workdir.clone())
            .stdout(std::process::Stdio::piped())
            .stderr(std::process::Stdio::piped())
            .spawn()
        {
            Ok(mut child) => {
                let stderr_reader = BufReader::new(child.stderr.take().unwrap());
                let stdout_reader = BufReader::new(child.stdout.take().unwrap());

                tokio::spawn(read_stream_and_capture(stderr_reader, app.id.clone(), true));
                tokio::spawn(read_stream_and_capture(
                    stdout_reader,
                    app.id.clone(),
                    false,
                ));

                self.app = Some(app.clone());
                self.program = Some(child);

                Ok(())
            }
            Err(err) => Err(err),
        };
    }

    pub async fn stop(&mut self) -> Result<(), io::Error> {
        if let Some(child) = self.program.borrow_mut() {
            return child.kill().await;
        }
        Ok(())
    }

    pub async fn get_stdout(&self) -> Option<String> {
        if let Some(app) = self.app.clone() {
            STDOUT.lock().await.get(&app.id).cloned()
        } else {
            None
        }
    }

    pub async fn get_stderr(&self) -> Option<String> {
        if let Some(app) = self.app.clone() {
            STDERR.lock().await.get(&app.id).cloned()
        } else {
            None
        }
    }
}

impl Default for AppInstance {
    fn default() -> Self {
        Self::new()
    }
}

async fn read_stream_and_capture<R>(reader: R, id: String, is_err: bool) -> io::Result<()>
where
    R: tokio::io::AsyncBufRead + Unpin,
{
    let mut lines = reader.lines();
    while let Some(line) = lines.next_line().await? {
        if !is_err {
            if let Some(out) = STDOUT.lock().await.get_mut(&id) {
                out.push_str(&line);
            }
        } else if let Some(out) = STDERR.lock().await.get_mut(&id) {
            out.push_str(&line);
        }
    }
    Ok(())
}
