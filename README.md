# AWS v4 Signing Proxy

Proxy to sign requests with AWS V4 Signatures. This project exists because I
got fed up with trying to sign requests using Java.

Is this a fantastic idea? No. Does it get the job done? Yes.

## Configuration

Configuration is performed via environment variables:

* ``AWS_SIGN_PROXY_BIND``: the address to bind to accept requests to sign and
  proxy to another service. Default is ``:8080``.
* ``AWS_SIGN_PROXY_TARGETPROTO``: the protocol to use when proxying the request
  to the other service. Default is ``https``.
* ``AWS_SIGN_PROXY_TARGETHOST``: hostname for the service which requires AWS v4
  signed requests.
* ``AWS_SIGN_PROXY_PROVIDER``: provider for the target service.
* ``AWS_SIGN_PROXY_REGION``: region for the target service. This will default
  to ``AWS_DEFAULT_REGION`` if not explicitly set.
* ``AWS_SIGN_PROXY_BLOCKHEADERS``: a list of header names that should be
  removed from the proxied request.
* ``AWS_SIGN_PROXY_EXTRAHEADERS``: a map of additional headers that should be
  added to the proxied request.

It is expected that you have ``AWS_ACCESS_KEY_ID`` and
``AWS_SECRET_ACCESS_KEY`` for the target service already set in your
environment as well.

## Usage

To issue a request to a backing service ``foo.bar.com`` using the ``baz``
provider in the ``us-east-1`` region:

```bash
$ export AWS_SIGN_PROXY_TARGETHOST=foo.bar.com
$ export AWS_SIGN_PROXY_PROVIDER=baz
$ export AWS_SIGN_PROXY_REGION=us-east-1
$ aws-sign-proxy
```

Now you may issue requests to the backing service by sending your plain
requests to ``http://localhost:8080/`` instead:

```bash
$ curl http://localhost:8080/some/service/endpoint -H 'X-Api-Key: blablabla'
```

Such a request would result in a signed request to

    https://foo.bar.com/some/service/endpoint

along with the ``X-Api-Key`` header (and any others specified on the plain
request).
