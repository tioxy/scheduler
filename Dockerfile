FROM golang:alpine
LABEL maintainer="Gabriel Tiossi"

ARG IN_BINARY="./scheduler"
ARG USER_ID=1001
ARG USER_NAME="scheduler"
ARG GROUP_NAME="scheduler"

USER root
RUN adduser -H -D -u $USER_ID $USER_NAME $GROUP_NAME 
COPY --chown=scheduler:scheduler $IN_BINARY /usr/bin/scheduler

USER scheduler
CMD ["/usr/bin/scheduler"]
