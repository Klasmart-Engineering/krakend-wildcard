# Microgateway Custom Plugin Template

This repository contains a template project for creating custom KrakenD plugins. _(To be injected into Microgateway deployments)_

Any microgateway that wants to leverage a plugin should initialize a git-submodule in their `/plugins` directory with the respective plugin

## Steps

1. Initialize the go module system

```sh
go mod init yourmodulename
```

2. To build your plugin run - this specifically builds it as a go plugin _(required by krakenD)_

```sh
go build -buildmode=plugin -o yourplugin.so .
```

3. Go into the `Makefile` and update the `build` target to have the same command as step 2. _(namely making the name of
   the outputted `.so` file is correct)_

## Make

The `Makefile` is required as it is used by other tooling (namely the `Dockerfile` for the microgateways) in order to
keep builds consistent across the company. As such please make sure the `all`, `build` and `clean` targets are all set
up to run correctly. Please feel free to add additional targets as required for your project, but please keep the 3
targets above limited in their scope.

## KrakenD Docs

[Intro to plugins](https://www.krakend.io/docs/extending/introduction/)

There are four different types of plugins you can write:

1. **[HTTP server plugins](https://www.krakend.io/docs/extending/http-server-plugins/)** (or handler plugins): They belong to the router layer and let you modify the request before KrakenD starts processing it, or modify the final response. You can have several plugins at once.
2. **[HTTP client plugins](https://www.krakend.io/docs/extending/http-client-plugins/)** (or proxy client plugins): They belong to the proxy layer and let you change how KrakenD interacts (as a client) with your backend services. You can have one plugin per backend.
3. **[Response Modifier plugins](https://www.krakend.io/docs/extending/plugin-modifiers/)**: They are strictly modifiers and let you change the responses received from your backends.
4. **[Request Modififer plugins](https://www.krakend.io/docs/extending/plugin-modifiers/)**: They are strictly modifiers and let you change the requests sent to your backends.
