--- note: this needs improvement, i'll get to it someday
create table if not exists app_logs (
  timestamp DateTime default now(),
  level LowCardinality(String),
  message Nullable(String),
  meta Nullable(String),
) engine = MergeTree()
order by (timestamp);
