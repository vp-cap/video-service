ARG SERIVCE_PATH="/go/src/vp-cap/video-service"

################## 1st Build Stage ####################
FROM golang:1.7.3 AS builder
LABEL stage=builder

WORKDIR $(SERIVCE_PATH)
ADD . .

ENV GO111MODULE=on

# Cache go mods based on go.sum/go.mod files
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o video-service

################## 2nd Build Stage ####################

FROM busybox:1-glibc

COPY --from=builder $(SERIVCE_PATH)/upload-service /usr/local/bin/video-service
COPY --from=builder $(SERIVCE_PATH)/config/config.yaml /usr/local/bin/config/config.yaml

ENTRYPOINT ["./usr/bin/video-service"]
