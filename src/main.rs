use axum::http::StatusCode;
use axum::routing::get;
use axum::Json;
use axum::Router;
use log::info;
use log::error;
use serde_json::json;
use serde_json::Value;
use std::env;
use std::process::exit;

fn send_standardized_response(
  code: StatusCode,
  message: &str,
  data: Value,
) -> (StatusCode, Json<Value>) {
  (
      code,
      Json(json!({
          "message": message,
          "data": data,
      })),
  )
}

#[tokio::main]
async fn main() {
    env_logger::Builder::new()
        .filter_level(log::LevelFilter::Info)
        .init();
    
    match dotenvy::dotenv() {
        Ok(_) => info!("dotenvy loaded"),
        Err(e) => {
          error!("dotenvy error: {}", e);
          exit(1);
        },
    }

    let port = env::var("PORT").unwrap_or("8084".to_string());
    let _api_key = env::var("NATTUKAKA_API_KEY").unwrap_or("NATTUKAKA_API_KEY_NOT_SET".to_string());

    let app: Router = Router::new()
        .route("/", get(hello));

    let listener = match tokio::net::TcpListener::bind(format!("{}:{}", "0.0.0.0", port))
    .await {
        Ok(listener) => listener,
        Err(e) => {
          error!("listener error: {}", e);
          exit(1);
        },
    };
    
    info!("server listening on port {}", port);
    match axum::serve(listener, app).await {
        Ok(_) => {},
        Err(e) => {
          error!("server error: {}", e);
          exit(1);
        },
    }
}

async fn hello() -> (StatusCode, Json<Value>) {
  send_standardized_response(StatusCode::OK, "nattukaka is up!", Value::Null)
}
