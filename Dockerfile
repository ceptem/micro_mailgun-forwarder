FROM        golang:latest AS build-image
LABEL       maintainer="raphael@ceptem.com"
RUN         go get gopkg.in/mailgun/mailgun-go.v1
ADD         src /src 
WORKDIR     /src 
ENV         CGO_ENABLED=0
ENV         GOOS=linux
RUN         go build -a -ldflags '-w -s' -installsuffix cgo -o app .
RUN         mkdir /1bsh58sd

FROM        scratch
COPY        --from=build-image /src/app /
COPY        --from=build-image /1bsh58sd/ /etc/ceptem/us/mailgun-forward-email/
VOLUME      /etc/ceptem/us/mailgun-forward-email/
EXPOSE      6543
ENTRYPOINT  ["/app"]
