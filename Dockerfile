FROM        golang:latest AS build-image
LABEL       maintainer="raphael@ceptem.com"
RUN         go get gopkg.in/mailgun/mailgun-go.v1
ADD         src /src 
WORKDIR     /src 
ENV         CGO_ENABLED=0
ENV         GOOS=linux
RUN         go build -a -ldflags '-w -s' -installsuffix cgo -o ../mailgun-forwarder .
RUN         mkdir /profiles

FROM        scratch
ENV         PROFILEDIR=/profiles
ENV         UID=0
ENV         GID=0
ENV         ADDRESS=0.0.0.0
ENV         PORT=6543
COPY        --from=build-image /mailgun-forwarder /profiles /
VOLUME      /profiles
EXPOSE      6543
CMD         ["/mailgun-forwarder"]
