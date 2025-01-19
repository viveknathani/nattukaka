use axum::body::Body;
use axum::http::Response;
use axum::http::StatusCode;
use axum::http::Uri;
use axum::response::Html;
use axum::response::IntoResponse;
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

async fn hello() -> (StatusCode, Html<String>) {
  let content = tokio::fs::read_to_string("static/index.html").await.unwrap();
  (StatusCode::OK, Html(content))
}

async fn serve_static(uri: Uri) -> Response<Body> {
  let static_dir = "static";
  let path = uri.path().trim_start_matches("/static/");
  let file_path = format!("{}/{}", static_dir, path);

  if let Ok(content) = tokio::fs::read(file_path).await {
    return (StatusCode::OK, content).into_response();
  } else {
    return (StatusCode::NOT_FOUND, "not found".to_string()).into_response();
  }
}

async fn handle_404() -> (StatusCode, Json<Value>) {
  send_standardized_response(
      StatusCode::NOT_FOUND,
      "not found",
      json!({}),
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
        .route("/", get(hello))
        .route("/static/{*wildcard}", get(serve_static))
        .fallback(handle_404);

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
