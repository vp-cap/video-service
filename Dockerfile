ARG SERVICE_PATH="/go/src/vp-cap/video-service"

################## 1st Build Stage ####################
FROM golang:1.15 AS builder
LABEL stage=builder
ARG SERVICE_PATH
ARG GIT_USER
ARG GIT_PASS

WORKDIR ${SERVICE_PATH}

ENV GO111MODULE=on
RUN git config --global url."https://$GIT_USER:$GIT_PASS@github.com".insteadOf "https://github.com"
RUN go env -w GOPRIVATE=github.com/vp-cap

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go install
# RUN ls

# ################## 2nd Build Stage ####################
FROM busybox:1-glibc
ARG SERVICE_PATH

COPY --from=builder /go/bin/video-service /usr/local/bin/video-service
COPY --from=builder ${SERVICE_PATH}/config.yaml /usr/local/bin/config.yaml
RUN cd /usr/local/bin && ls

ENTRYPOINT ["./usr/local/bin/video-service"]