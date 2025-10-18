FROM golang:1.23-alpine AS builder
WORKDIR /workspace
RUN apk add --no-cache git
COPY . .
ENV CGO_ENABLED=0

RUN set -eux; \
    for dir in TrimVideo EditVideo AudioRemoval Watermarking gossipOpenClose StatesMachine AdminCache; do \
      cd /workspace/$dir && go mod download; \
    done

RUN set -eux; \
    mkdir -p /out/trim-video /out/edit-video /out/audio-removal /out/watermarking /out/gossip-open-close /out/states-machine /out/admin-cache; \
    cd /workspace/TrimVideo && go build -o /out/trim-video/trimvideo ./cmd/main.go; \
    cd /workspace/EditVideo && go build -o /out/edit-video/editvideo ./cmd/main.go; \
    cd /workspace/AudioRemoval && go build -o /out/audio-removal/audioremoval ./cmd/main.go; \
    cd /workspace/Watermarking && go build -o /out/watermarking/watermarking ./cmd/main.go && cp -r assets /out/watermarking/; \
    cd /workspace/gossipOpenClose && go build -o /out/gossip-open-close/gossipopenclose ./cmd/main.go && cp -r assets /out/gossip-open-close/; \
    cd /workspace/StatesMachine && go build -o /out/states-machine/statesmachine ./cmd/main.go; \
    cd /workspace/AdminCache && go build -o /out/admin-cache/admincache ./cmd/admincache; \
    true

FROM alpine:3.20
RUN apk add --no-cache ffmpeg ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/ /app/
COPY start-workers.sh /usr/local/bin/start-workers.sh
RUN chmod +x /usr/local/bin/start-workers.sh \
    && addgroup -S workers \
    && adduser -S worker -G workers \
    && chown -R worker:workers /app \
    && mkdir -p /tmp \
    && chmod 1777 /tmp

ENV RABBITMQ_URL=amqp://admin:admin@10.0.5.108:5672 \
    AWS_REGION=us-east-1 \
    AWS_ACCESS_KEY_ID=changeme \
    AWS_SECRET_ACCESS_KEY=changeme \
    S3_ENDPOINT= \
    S3_USE_PATH_STYLE=false \
    STATE_MACHINE_QUEUE=states_machine_queue \
    EDIT_VIDEO_QUEUE=edit_video_queue \
    AUDIO_REMOVAL_QUEUE=audio_removal_queue \
    WATERMARKING_QUEUE=watermarking_queue \
    DATABASE_URL=postgres://app_user:app_password@videorank.c1o2s0umijp6.us-east-1.rds.amazonaws.com:5432/videorank?sslmode=require \
    REDIS_ADDR=10.0.3.78:6379 \
    TRIM_VIDEO_QUEUE_NAME=trim_video_queue \
    TRIM_VIDEO_RAW_BUCKET=raw-videos \
    TRIM_VIDEO_PROCESSED_BUCKET=processed-videos-trim \
    TRIM_VIDEO_MAX_SECONDS=30 \
    TRIM_VIDEO_MAX_RETRIES=3 \
    EDIT_VIDEO_QUEUE_NAME=edit_video_queue \
    EDIT_VIDEO_RAW_BUCKET=processed-videos-trim \
    EDIT_VIDEO_PROCESSED_BUCKET=processed-videos-edit \
    EDIT_VIDEO_MAX_SECONDS=30 \
    EDIT_VIDEO_MAX_RETRIES=3 \
    AUDIO_REMOVAL_QUEUE_NAME=audio_removal_queue \
    AUDIO_REMOVAL_RAW_BUCKET=processed-videos-edit \
    AUDIO_REMOVAL_PROCESSED_BUCKET=processed-videos-audio-removal \
    AUDIO_REMOVAL_MAX_RETRIES=3 \
    WATERMARKING_QUEUE_NAME=watermarking_queue \
    WATERMARKING_RAW_BUCKET=processed-videos-audio-removal \
    WATERMARKING_PROCESSED_BUCKET=processed-videos-watermarking \
    WATERMARKING_MAX_SECONDS=30 \
    WATERMARKING_MAX_RETRIES=3 \
    GOSSIP_QUEUE_NAME=gossip_open_close_queue \
    GOSSIP_S3_BUCKET_RAW=processed-videos-watermarking \
    GOSSIP_S3_BUCKET_PROCESSED=processed-videos \
    GOSSIP_MAX_SECONDS=30 \
    GOSSIP_INTRO_SECONDS=2.5 \
    GOSSIP_OUTRO_SECONDS=2.5 \
    GOSSIP_TARGET_WIDTH=1280 \
    GOSSIP_TARGET_HEIGHT=720 \
    GOSSIP_FPS=30 \
    GOSSIP_MAX_RETRIES=3 \
    STATES_MACHINE_MAX_RETRIES=3 \
    STATES_MACHINE_RETRY_DELAY_MINUTES=5 \
    ADMIN_CACHE_REDIS_ADDR=10.0.3.78:6379 \
    ADMIN_CACHE_POSTGRES_DSN=postgres://app_user:app_password@videorank.c1o2s0umijp6.us-east-1.rds.amazonaws.com:5432/videorank?sslmode=require  \
    REFRESH_INTERVAL_SECONDS=30

USER worker
ENTRYPOINT ["/usr/local/bin/start-workers.sh"]
