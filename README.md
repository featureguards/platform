# Platform

This repository contains all the backend/frontend code for [featureguards.com](https://www.featureguards.com). This was a fun project to build from scratch. The outcome is a
pretty useful platform for building modern webservices.

A series of posts will go over the exciting technical aspects of building this service. Please, follow me at [Infrastructure tales](https://brawi.substack.com).

# How to use?

This project was a great experience building something completely from scratch.

# Code Structure

```
/-
  |---- .air -> directory for running 'live-reloadable' Go binaries.
  |---- .envs -> directory for having .env files for various environments.
  |---- .vscode -> project-wide vs code settings
  |---- app -> dashboard frontend code
  |---- certs -> local development certs used for Envoy to serve over TLS (similar to production).
  |---- envoy -> Envoy configuration files
  |---- go -> backend Go code (Dashboard, Auth, and Toggles services)
  |---- openapi -> openapi generated APIs from gRPC proto definitions. Mostly used for browser javascript.
  |---- ory/kratos -> Ory Kratos configs for user auth/management.
  |---- proto -> protobuf and gRPC definitions and code generation scripts.

```

# Setting Up

1. Create a local directory `featureguards`
2. `cd featureguards`
3. Clone repo via `git clone --depth 1 git@github.com:featureguards/platform.git`
4. **You must clone other repos too since protobuf generation uses local directories.**
5. Clone the following:
   - `git clone --depth 1 git@github.com:featureguards/featureguards-go.git`
   - `git clone --depth 1 git@github.com:featureguards/featureguards-js.git`
   - `git clone --depth 1 git@github.com:featureguards/featureguards-py.git`.
6. Install docker, node, npm, mutagen, protobuf, Go, and Python and include relevant binaries in path.
7. Add the following to `/etc/hosts`

```
 127.0.0.1 app.featureguards.dev
 ::1       app.featureguards.dev
 127.0.0.1 api.featureguards.dev
 ::1       api.featureguards.dev
```

8. Trust the local certs by importing them into Key Chain and trusting them so Chrome can the local TLS service.
9. Add the following to `~/.zshrc` or `~/.bashrc` depending on your shell

```
 alias dd='COMPOSE_PROJECT_NAME=platform mutagen-compose -f local.yml'
224 alias dt='COMPOSE_PROJECT_NAME=platform-test mutagen-compose -f test.yml'
```

# Running Locally

To run the application locally after reloading `.zshrc` or `.bashrc` simply type

```
dd up
```

This will bring up all containers based on `local.yml` file config.

If everything is running, use Chrome and go to `https://app.featureguards.dev`. It should load the web application.

If there are errors, everything should be written on the commandline output. You can use it for debugging.

<b> NOTE: <br>
Never build the packages (node or Go) locally on your mac. This is because they need to be built inside the container to match Linux environment <b>

# Testing

There are some Go unit-tests and integration tests. They mostly use mocks for redis/sqlite. To run the tests, you need to type the following

```
dt up

```

# FAQ

1. Why use a custom domain and certificates locally?

   > This makes testing very similar to production. No hadcoded `localhost` or HTTP vs HTTPs anywhere.

2. I updated the gRPC services. How do I generate the protobuf stubs?

   > `cd proto; make`

3. How do I get `Sign in With Google` to work?

   > You need to include your own `client_id` and `client_secret` in `.envs/local/kratos`.

4. How to host this in the cloud?
   > There are various options. My setup is on AWS with the following services:

- 2 ALBs: One for serving app frontend code and dashboard based on the route and the other ALB is for serving gRPC/HTTP APIs for client SDKs.
- ECS: 3 services (App, Api, and Dashboard). Docker files are included in the repo as aws.Dockerfile. The API service currently includes both toggles and auth services. You can split services horizontally.
- Aurora Postgres RDS.
- ElastiCache Redis cluster.
- ECR: All images are built locally and served from ECR to avoid dependency changes and bundle some configs.
- Codebuild: 4 images are built: Kratos, typescript, Go, and Envoy.
- Parameter Store: All secrets are stored in Systems Manager and are loaded upon spinning up the ECS task.
