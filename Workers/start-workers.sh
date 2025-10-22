#!/bin/sh
set -eu

log() {
  printf '%s %s\n' "$(date '+%Y-%m-%dT%H:%M:%S%z')" "$1"
}

require_env() {
  var="$1"
  eval "value=\${$var:-}"
  if [ -z "$value" ]; then
    log "Environment variable $var is required"
    exit 1
  fi
}

PIDS=""
PID_NAMES=""

track_pid() {
  pid="$1"
  name="$2"
  PIDS="$PIDS $pid"
  PID_NAMES="$PID_NAMES $pid:$name"
}

get_name_by_pid() {
  target="$1"
  for entry in $PID_NAMES; do
    case "$entry" in
      "$target":*)
        echo "${entry#*:}"
        return 0
        ;;
    esac
  done
  echo "unknown"
}

stop_all() {
  for pid in $PIDS; do
    if kill -0 "$pid" 2>/dev/null; then
      name="$(get_name_by_pid "$pid")"
      log "Stopping $name (pid $pid)"
      kill "$pid" 2>/dev/null || true
    fi
  done
}

terminate() {
  log "Termination signal received; shutting down workers"
  stop_all
  wait || true
  exit 0
}

trap terminate INT TERM


# Shared defaults
RABBITMQ_URL="${RABBITMQ_URL:-amqp://admin:admin@rabbitmq:5672/}"
STATE_MACHINE_QUEUE="${STATE_MACHINE_QUEUE:-states_machine_queue}"
EDIT_VIDEO_QUEUE="${EDIT_VIDEO_QUEUE:-edit_video_queue}"
AUDIO_REMOVAL_QUEUE="${AUDIO_REMOVAL_QUEUE:-audio_removal_queue}"
WATERMARKING_QUEUE="${WATERMARKING_QUEUE:-watermarking_queue}"
DATABASE_URL="${DATABASE_URL:-postgres://app_user:app_password@postgres:5432/videorank?sslmode=disable}"
REDIS_ADDR="${REDIS_ADDR:-redis:6379}"
require_env AWS_REGION
require_env AWS_ACCESS_KEY_ID
require_env AWS_SECRET_ACCESS_KEY
S3_ENDPOINT="${S3_ENDPOINT:-}"
S3_USE_PATH_STYLE="${S3_USE_PATH_STYLE:-false}"


start_trim_video() {
  require_env TRIM_VIDEO_QUEUE_NAME
  require_env TRIM_VIDEO_RAW_BUCKET
  require_env TRIM_VIDEO_PROCESSED_BUCKET
  (
    cd /app/trim-video
    exec env \
      RABBITMQ_URL="$RABBITMQ_URL" \
      AWS_REGION="$AWS_REGION" \
      AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
      AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
      S3_ENDPOINT="$S3_ENDPOINT" \
      S3_USE_PATH_STYLE="$S3_USE_PATH_STYLE" \
      RAW_BUCKET="$TRIM_VIDEO_RAW_BUCKET" \
      PROCESSED_BUCKET="$TRIM_VIDEO_PROCESSED_BUCKET" \
      QUEUE_NAME="$TRIM_VIDEO_QUEUE_NAME" \
      STATE_MACHINE_QUEUE="$STATE_MACHINE_QUEUE" \
      MAX_SECONDS="${TRIM_VIDEO_MAX_SECONDS:-30}" \
      MAX_RETRIES="${TRIM_VIDEO_MAX_RETRIES:-3}" \
      QUEUE_MAX_LENGTH="${TRIM_VIDEO_QUEUE_MAX_LENGTH:-1000}" \
      ./trimvideo
  ) &
  track_pid "$!" "trim-video"
}

start_edit_video() {
  require_env EDIT_VIDEO_QUEUE_NAME
  require_env EDIT_VIDEO_RAW_BUCKET
  require_env EDIT_VIDEO_PROCESSED_BUCKET
  (
    cd /app/edit-video
    exec env \
      RABBITMQ_URL="$RABBITMQ_URL" \
      AWS_REGION="$AWS_REGION" \
      AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
      AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
      S3_ENDPOINT="$S3_ENDPOINT" \
      S3_USE_PATH_STYLE="$S3_USE_PATH_STYLE" \
      RAW_BUCKET="$EDIT_VIDEO_RAW_BUCKET" \
      PROCESSED_BUCKET="$EDIT_VIDEO_PROCESSED_BUCKET" \
      QUEUE_NAME="$EDIT_VIDEO_QUEUE_NAME" \
      STATE_MACHINE_QUEUE="$STATE_MACHINE_QUEUE" \
      MAX_SECONDS="${EDIT_VIDEO_MAX_SECONDS:-30}" \
      MAX_RETRIES="${EDIT_VIDEO_MAX_RETRIES:-3}" \
      QUEUE_MAX_LENGTH="${EDIT_VIDEO_QUEUE_MAX_LENGTH:-1000}" \
      ./editvideo
  ) &
  track_pid "$!" "edit-video"
}

start_audio_removal() {
  require_env AUDIO_REMOVAL_QUEUE_NAME
  require_env AUDIO_REMOVAL_RAW_BUCKET
  require_env AUDIO_REMOVAL_PROCESSED_BUCKET
  (
    cd /app/audio-removal
    exec env \
      RABBITMQ_URL="$RABBITMQ_URL" \
      AWS_REGION="$AWS_REGION" \
      AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
      AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
      S3_ENDPOINT="$S3_ENDPOINT" \
      S3_USE_PATH_STYLE="$S3_USE_PATH_STYLE" \
      RAW_BUCKET="$AUDIO_REMOVAL_RAW_BUCKET" \
      PROCESSED_BUCKET="$AUDIO_REMOVAL_PROCESSED_BUCKET" \
      QUEUE_NAME="$AUDIO_REMOVAL_QUEUE_NAME" \
      STATE_MACHINE_QUEUE="$STATE_MACHINE_QUEUE" \
      MAX_RETRIES="${AUDIO_REMOVAL_MAX_RETRIES:-3}" \
      QUEUE_MAX_LENGTH="${AUDIO_REMOVAL_QUEUE_MAX_LENGTH:-1000}" \
      ./audioremoval
  ) &
  track_pid "$!" "audio-removal"
}

start_watermarking() {
  require_env WATERMARKING_QUEUE_NAME
  require_env WATERMARKING_RAW_BUCKET
  require_env WATERMARKING_PROCESSED_BUCKET
  (
    cd /app/watermarking
    exec env \
      RABBITMQ_URL="$RABBITMQ_URL" \
      AWS_REGION="$AWS_REGION" \
      AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
      AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
      S3_ENDPOINT="$S3_ENDPOINT" \
      S3_USE_PATH_STYLE="$S3_USE_PATH_STYLE" \
      RAW_BUCKET="$WATERMARKING_RAW_BUCKET" \
      PROCESSED_BUCKET="$WATERMARKING_PROCESSED_BUCKET" \
      QUEUE_NAME="$WATERMARKING_QUEUE_NAME" \
      STATE_MACHINE_QUEUE="$STATE_MACHINE_QUEUE" \
      MAX_SECONDS="${WATERMARKING_MAX_SECONDS:-30}" \
      MAX_RETRIES="${WATERMARKING_MAX_RETRIES:-3}" \
      QUEUE_MAX_LENGTH="${WATERMARKING_QUEUE_MAX_LENGTH:-1000}" \
      ./watermarking
  ) &
  track_pid "$!" "watermarking"
}

start_gossip_open_close() {
  require_env GOSSIP_QUEUE_NAME
  require_env GOSSIP_S3_BUCKET_RAW
  require_env GOSSIP_S3_BUCKET_PROCESSED
  (
    cd /app/gossip-open-close
    exec env \
      RABBITMQ_URL="$RABBITMQ_URL" \
      AWS_REGION="$AWS_REGION" \
      AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
      AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
      S3_ENDPOINT="$S3_ENDPOINT" \
      S3_USE_PATH_STYLE="$S3_USE_PATH_STYLE" \
      S3_BUCKET_RAW="$GOSSIP_S3_BUCKET_RAW" \
      S3_BUCKET_PROCESSED="$GOSSIP_S3_BUCKET_PROCESSED" \
      QUEUE_NAME="$GOSSIP_QUEUE_NAME" \
      MAX_SECONDS="${GOSSIP_MAX_SECONDS:-30}" \
      INTRO_SECONDS="${GOSSIP_INTRO_SECONDS:-2.5}" \
      OUTRO_SECONDS="${GOSSIP_OUTRO_SECONDS:-2.5}" \
      TARGET_WIDTH="${GOSSIP_TARGET_WIDTH:-1280}" \
      TARGET_HEIGHT="${GOSSIP_TARGET_HEIGHT:-720}" \
      FPS="${GOSSIP_FPS:-30}" \
      LOGO_PATH="${GOSSIP_LOGO_PATH:-/app/gossip-open-close/assets/nba-logo-removebg-preview.png}" \
      MAX_RETRIES="${GOSSIP_MAX_RETRIES:-3}" \
      QUEUE_MAX_LENGTH="${GOSSIP_QUEUE_MAX_LENGTH:-1000}" \
      ./gossipopenclose
  ) &
  track_pid "$!" "gossip-open-close"
}

start_states_machine() {
  (
    cd /app/states-machine
    exec env \
      RABBITMQ_URL="$RABBITMQ_URL" \
      QUEUE_NAME="$STATE_MACHINE_QUEUE" \
      EDIT_VIDEO_QUEUE="$EDIT_VIDEO_QUEUE" \
      AUDIO_REMOVAL_QUEUE="$AUDIO_REMOVAL_QUEUE" \
      WATERMARKING_QUEUE="$WATERMARKING_QUEUE" \
      DATABASE_URL="${STATES_MACHINE_DATABASE_URL:-$DATABASE_URL}" \
      MAX_RETRIES="${STATES_MACHINE_MAX_RETRIES:-3}" \
      RETRY_DELAY_MINUTES="${STATES_MACHINE_RETRY_DELAY_MINUTES:-5}" \
      ./statesmachine
  ) &
  track_pid "$!" "states-machine"
}

start_admin_cache() {
  (
    cd /app/admin-cache
    exec env \
      REDIS_ADDR="${ADMIN_CACHE_REDIS_ADDR:-$REDIS_ADDR}" \
      CACHE_PREFIX="${ADMIN_CACHE_CACHE_PREFIX:-videorank:}" \
      SCHEMA_VERSION="${ADMIN_CACHE_SCHEMA_VERSION:-v2}" \
      CACHE_TTL_FRESH_SECONDS="${ADMIN_CACHE_CACHE_TTL_FRESH_SECONDS:-900}" \
      CACHE_MAX_STALE_SECONDS="${ADMIN_CACHE_CACHE_MAX_STALE_SECONDS:-600}" \
      CACHE_LOCK_LEASE_SECONDS="${ADMIN_CACHE_CACHE_LOCK_LEASE_SECONDS:-10}" \
      CACHE_JITTER_PERCENT="${ADMIN_CACHE_CACHE_JITTER_PERCENT:-10}" \
      POSTGRES_DSN="${ADMIN_CACHE_POSTGRES_DSN:-$DATABASE_URL}" \
      DB_READ_TIMEOUT_SECONDS="${ADMIN_CACHE_DB_READ_TIMEOUT_SECONDS:-3}" \
      DB_MAX_RETRIES="${ADMIN_CACHE_DB_MAX_RETRIES:-3}" \
      REFRESH_INTERVAL_SECONDS="${ADMIN_CACHE_REFRESH_INTERVAL_SECONDS:-300}" \
      BATCH_SIZE_CITIES="${ADMIN_CACHE_BATCH_SIZE_CITIES:-50}" \
      CACHE_MAX_TOP_USERS="${ADMIN_CACHE_CACHE_MAX_TOP_USERS:-10}" \
      WARM_CITIES="${ADMIN_CACHE_WARM_CITIES:-}" \
      ./admincache
  ) &
  track_pid "$!" "admin-cache"
}

log "Starting workers without waiting for external services"

start_trim_video
start_edit_video
start_audio_removal
start_watermarking
start_gossip_open_close
start_admin_cache
start_states_machine

log "All workers started"

set +e
status=0
for pid in $PIDS; do
  if wait "$pid"; then
    log "Process $(get_name_by_pid "$pid") (pid $pid) exited normally"
  else
    rc=$?
    name="$(get_name_by_pid "$pid")"
    log "Process $name (pid $pid) exited with status $rc"
    status=$rc
    stop_all
    break
  fi
done
wait || true
exit "$status"
