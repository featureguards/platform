# syntax=docker/dockerfile:1.4
FROM envoyproxy/envoy:v1.21-latest

EXPOSE 443
EXPOSE 8001
EXPOSE 80

COPY prod/envoy.yaml /etc/envoy/
COPY ./dashboard.pb /protos
COPY ./auth.pb /protos
COPY ./toggles.pb /protos