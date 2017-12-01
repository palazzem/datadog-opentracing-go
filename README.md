# Datadog Example

This repository contains an example Go app that is instrumented using the new
[OpenTracing API][1] available in the [Datadog Go Client][2].

## Run the example

To run the service, simply::

    $ go run main.go

In another shell, you can simply call your service::

    $ curl localhost:3000/account/42

[1]: https://github.com/opentracing/opentracing-go
[2]: https://github.com/DataDog/dd-trace-go
