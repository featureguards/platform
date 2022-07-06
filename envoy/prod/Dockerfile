# syntax=docker/dockerfile:1.4
FROM 351563604431.dkr.ecr.us-west-2.amazonaws.com/envoyproxy/envoy:v1.21-latest

EXPOSE 8443
EXPOSE 8001
EXPOSE 8080

COPY prod/envoy.yaml /etc/envoy/
COPY ./dashboard.pb /protos/dashboard.pb
COPY ./auth.pb /protos/auth.pb
COPY ./toggles.pb /protos/toggles.pb