# Traefik Plugin: staticresponse

This repository contains a Traefik middleware used for hijacking requests and responsing with a predefined code, headers and body. The request is never sent to the backend.

Heavily inspired by the `plugindemo` and `noop` traefik plugins.

## Usage

Follow traefik's instructions https://plugins.traefik.io/install.

```yml
experimental:
  plugins:
    staticresponse:
      moduleName: "github.com/jdel/staticresponse"
      version: "v0.0.1"
```

## Configuration

Check the `docker-compose.yml` file for examples with docker labels.

Check the `example-dynamic-config.yml` file for examples with config file.
