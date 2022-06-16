# Proto Generation

### Envoy Descriptor

See [Envoy](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/grpc_json_transcoder_filter) for more information

```
git clone https://github.com/googleapis/googleapis
```

```
protoc --include_imports --include_source_info \
    --descriptor_set_out=proto.pb *.proto
```

### Golang Generation

```
protoc --go_out=../go --go_opt=module=platform/go \
    --go-grpc_out=../go --go-grpc_opt=module=platform/go *.proto
```

### OpenAPI Generation

```
protoc *.proto --openapi_out=../openapi
```

### Typescript Generation

```
openapi-generator generate -i openapi.yaml -g typescript-axios -o ../app/api
```
