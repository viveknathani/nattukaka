use axum::body::Body;
use axum::http::Response;
use axum::http::StatusCode;
use axum::http::Uri;
use axum::response::Html;
use axum::response::IntoResponse;
use axum::routing::get;
use axum::Json;
use axum::Router;
use chrono::Utc;
use log::error;
use log::info;
use redis::Commands;
use redis::Connection;
use serde::Deserialize;
use serde::Deserializer;
use serde::Serialize;
use serde::Serializer;
use serde_json::json;
use serde_json::Value;
use std::env;
use std::process::exit;

const SERVICE_COUNT: usize = 3;
const SERVICES: [&str; SERVICE_COUNT] = ["teachyourselfmath", "vivekn.dev", "workdiff"];
const REDIS_KEYS_PREFIX_NATTUKAKA_SERVICE: &str = "NATTUKAKA:DEPLOYMENTS";

#[derive(Serialize, Deserialize, Debug)]
enum DeploymentStatus {
    NotAvailable,
    InProgress,
    Success,
    Failed,
}

fn serialize_deployment_status<S>(
    status: &DeploymentStatus,
    serializer: S,
) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    serializer.serialize_str(&status.to_string())
}

fn deserialize_deployment_status<'de, D>(deserializer: D) -> Result<DeploymentStatus, D::Error>
where
    D: Deserializer<'de>,
{
    let s = String::deserialize(deserializer)?;
    match DeploymentStatus::from_str(&s) {
        Some(status) => Ok(status),
        None => Err(serde::de::Error::custom(format!("Invalid deployment status: {}", s))),
    }
}

impl DeploymentStatus {
    fn to_string(&self) -> &str {
        match self {
            DeploymentStatus::NotAvailable => "NOT_AVAILABLE",
            DeploymentStatus::InProgress => "IN_PROGRESS",
            DeploymentStatus::Success => "SUCCESS",
            DeploymentStatus::Failed => "FAILED",
        }
    }

    fn from_str(s: &str) -> Option<DeploymentStatus> {
      match s {
          "NOT_AVAILABLE" => Some(DeploymentStatus::NotAvailable),
          "IN_PROGRESS" => Some(DeploymentStatus::InProgress),
          "SUCCESS" => Some(DeploymentStatus::Success),
          "FAILED" => Some(DeploymentStatus::Failed),
          _ => None,
      }
  }
}

#[derive(Serialize, Deserialize, Debug)]
struct Deployment {
    service: String,

    #[serde(serialize_with = "serialize_deployment_status", deserialize_with = "deserialize_deployment_status")]
    status: DeploymentStatus,

    timestamp: String,
}

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
    let content = tokio::fs::read_to_string("static/index.html")
        .await
        .unwrap();
    (StatusCode::OK, Html(content))
}

async fn serve_static(uri: Uri) -> Response<Body> {
    let static_dir = "static";
    let path = uri.path().trim_start_matches("/static/");
    let file_path = format!("{}/{}", static_dir, path);

    if let Ok(content) = tokio::fs::read(file_path).await {
        return (StatusCode::OK, content).into_response();
    } else {
        return (StatusCode::NOT_FOUND, "are you lost?".to_string()).into_response();
    }
}

async fn handle_404() -> (StatusCode, Json<Value>) {
    send_standardized_response(StatusCode::NOT_FOUND, "are you lost?", json!({}))
}

async fn fetch_deployments(redis_conn: &mut Connection) -> Vec<Deployment> {
    let mut result: Vec<Deployment> = Vec::new();

    for service in SERVICES {
        let redis_key = format!("{}:{}", REDIS_KEYS_PREFIX_NATTUKAKA_SERVICE, service);
        let out: Option<String> = match redis_conn.get(&redis_key) {
            Ok(val) => val,
            Err(err) => {
                error!(
                    "failed to get service info for {} with error: {}",
                    service, err
                );
                break;
            }
        };
        if out.is_some() {
            let deployment: Deployment = match serde_json::from_str(out.unwrap().as_str()) {
                Ok(val) => val,
                Err(err) => {
                    error!(
                        "failed to get service info for {} with error: {}",
                        service, err
                    );
                    break;
                }
            };
            result.push(deployment);
        }
    }

    return result;
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
        }
    }

    let port = env::var("PORT").unwrap_or_else(|_| "8084".to_string());

    let _api_key = env::var("NATTUKAKA_API_KEY").unwrap_or("NATTUKAKA_API_KEY_NOT_SET".to_string());

    let redis_url = env::var("REDIS_URL").unwrap_or("REDIS_URL_NOT_SET".to_string());

    let redis_client = match redis::Client::open(redis_url) {
        Ok(client) => client,
        Err(err) => {
            error!("could not get a client for redis: {}", err);
            exit(1);
        }
    };

    let redis_pool = match r2d2::Pool::builder().build(redis_client) {
        Ok(pool) => pool,
        Err(err) => {
            error!("could not get a pool for redis: {}", err);
            exit(1);
        }
    };

    let mut redis_conn = match redis_pool.get() {
        Ok(conn) => conn,
        Err(err) => {
            error!("could not connect with redis: {}", err);
            exit(1);
        }
    };

    for service in SERVICES {
        let redis_key = format!("{}:{}", REDIS_KEYS_PREFIX_NATTUKAKA_SERVICE, service);
        let out: Option<String> = match redis_conn.get(&redis_key) {
            Ok(val) => val,
            Err(err) => {
                error!(
                    "failed to get service info for {} with error: {}",
                    service, err
                );
                exit(1);
            }
        };
        if out.is_none() {
            let deployment = Deployment {
                service: (&service).to_string(),
                status: DeploymentStatus::NotAvailable,
                timestamp: Utc::now().to_rfc3339(),
            };

            let deployment_str = match serde_json::to_string(&deployment) {
                Ok(result) => result,
                Err(err) => {
                    error!(
                        "failed to serialize for {:?} with error: {}",
                        deployment, err
                    );
                    exit(1);
                }
            };

            let _: Option<String> = match redis_conn.set(&redis_key, deployment_str) {
                Ok(val) => val,
                Err(err) => {
                    error!("redis set for {} failed {}", service, err);
                    exit(1);
                }
            };
        }
        info!("service info for {} is set", service);
    }

    let redis_pool = redis_pool.clone();
    let app: Router = Router::new()
        .route("/", get(hello))
        .route("/static/{*wildcard}", get(serve_static))
        .route(
            "/deployments",
            get(|| async move {
                let mut connection = redis_pool.get().unwrap();
                let deployments = fetch_deployments(&mut connection).await;

                if deployments.len() < SERVICE_COUNT {
                    send_standardized_response(
                        StatusCode::INTERNAL_SERVER_ERROR,
                        "this should not happen, we don't have all the data :(",
                        json!(""),
                    )
                } else {
                    send_standardized_response(StatusCode::OK, "", json!(deployments))
                }
            }),
        )
        .fallback(handle_404);

    let listener = match tokio::net::TcpListener::bind(format!("{}:{}", "0.0.0.0", port)).await {
        Ok(listener) => listener,
        Err(e) => {
            error!("listener error: {}", e);
            exit(1);
        }
    };

    info!("server listening on port {}", port);
    match axum::serve(listener, app).await {
        Ok(_) => {}
        Err(e) => {
            error!("server error: {}", e);
            exit(1);
        }
    }
}
