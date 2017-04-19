Swagger2Postman in Go
=====================

[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

`swagger2postman` creates a Postman 2.0 Collection spec from a Swagger (OAI) 2.0 spec.

# Features

* Build on a base Postman 2.0 spec that contains information not readily stored in Swagger Spec to support Postman features such as scripts.
* Supports override parameters such as a Postman environment parameter for the URL hostname, e.g. https://{{MY_HOSTNAME}}/rest
* Supports additional headers, e.g. authorization headers using enviroment variables, e.g. `Authorization: Bearer {{myAccessToken}}`

# Notes

* Postman 4.10.7 does not natively support JSON requests so request bodies need to be entered using the raw body editor. A future task is to add Swagger request examples as default Postman request bodies.
* Postman 2.0 spec doesn't have a canonical schema. The `request.url` property can be populated by a URL string or a URL object. Swagger2Postman uses the URL object since it is more flexible. The function `simple.NewCanonicalCollectionFromBytes(bytes)` can be used to read either a simple or object based spec into a canonical object spec.

# Example Usage

There is an example conversion available, [`examples/ringcentral/convert.go`](https://github.com/grokify/swagger2postman-go/blob/master/examples/ringcentral/convert.go) which creates a Postman 2.0 spec for the [RingCentral REST API](https://developers.ringcentral.com) using a base Postman 2.0 spec and the RingCentral basic Swagger 2.0 spec.

Example files include:

* [RingCentral Swagger 2.0 spec](https://github.com/grokify/swagger2postman-go/blob/master/examples/ringcentral/ringcentral.swagger2.basic.json)
* [RingCentral Postman 2.0 base](https://github.com/grokify/swagger2postman-go/blob/master/examples/ringcentral/ringcentral.postman2.base.json)
* [RingCentral Postman 2.0 spec](https://github.com/grokify/swagger2postman-go/blob/master/examples/ringcentral/ringcentral.postman2.basic.json) - Import this into Postman

The RingCentral spec uses the following environment variables. The following is the Postman bulk edit format:

```
RC_SERVER_HOSTNAME:platform.devtest.ringcentral.com
RC_APP_KEY:myAppKey
RC_APP_SECRET:myAppSecret
RC_USER_USERNAME:myMainCompanyPhoneNumber
RC_USER_EXTENSION:myExtension
RC_USER_PASSWORD:myPassword
```

For multiple apps or users, simply create a different Postman environment for each.

To set your environment variables, use the Settings Gear icon and then click "Manage Environments"

 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/swagger2postman-go
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/swagger2postman-go
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/swagger2postman-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/swagger2postman-go/blob/master/LICENSE.md
