FROM alpine:3.8 AS build
RUN  apk add --no-cache ca-certificates tzdata
ADD build/. /.

FROM        scratch
COPY        --from=build /etc/ssl/certs /etc/ssl/certs
COPY        --from=build /usr/share/zoneinfo/Australia /usr/share/zoneinfo/Australia
COPY        --from=build /api /api
EXPOSE      8080
ENTRYPOINT  ["/api"]
