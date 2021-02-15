Swaggman - OpenAPI / Swagger Spec Utility and Postman Converter
===============================================================

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

![](docs/images/logo_swaggman_600x150.png "")

Swaggman is a multi-purpose OpenAPI and Postman utility.

**Major Features**

* CLI and library to Convert OpenAPI Specs to Postman Collection
  * Add Postman environment variables to URLs, e.g. Server URLs like `https://{{HOSTNAME}}/restapi`
  * Add headers, such as environment variable based Authorization headers, such as `Authorization: Bearer {{myAccessToken}}`
  * Utilize baseline Postman collection to add Postman-specific functionality including Postman `prerequest` scripts.
  * Add example request bodies, e.g. JSON bodies with example parameter values.
* OpenAPI 3 Spec SDK
  * Merge Multiple OAS3 Specs
  * Validate OAS3 Specs
  * Programmatically examine and modify OAS3 Specs
  * [Programmatically auto-correct OAS3 Specs](docs/openapi3_fix.md)
* OpenAPI 2 Spec SDK
  * Convert OpenAPI v2 to Postman Collection v2
  * Merge Multiple OAS2 Specs

# Features

These are all used in the included example discussed below.

[Additional discussion is available on Medium](https://medium.com/ringcentral-developers/using-postman-with-swagger-and-the-ringcentral-api-523712f792a0).

# Notes

* Postman 4.10.7 does not natively support JSON requests so request bodies need to be entered using the raw body editor. A future task is to add Swagger request examples as default Postman request bodies.
* Postman 2.0 spec supports polymorphism and doesn't have a canonical schema. For example, the `request.url` property can be populated by a URL string or a URL object. Swaggman uses the URL object since it is more flexible. The function `simple.NewCanonicalCollectionFromBytes(bytes)` can be used to read either a simple or object based spec into a canonical object spec.
* This has only been used on the RingCentral Swagger spec to date but will be used for more in the future. Please feel free to use and contribute. Examples are located in the `examples` folder.

# Installation

The following command will install the executable binary `swaggman` into the `~/go/bin` directory.

```bash
$ go get github.com/grokify/swaggman
```

# Usage

## Simple Usage

```
// Instantiate a converter with default configuration
conv := swaggman.NewConverter(swaggman.Configuration{})

// Convert a Swagger spec
err := conv.Convert("path/to/swagger.json", "path/to/pman.out.json")
```

## Usage with Features

The following can be added which are especially useful to use with environment variables.

* Custom Hostname
* Custom Headers

```
// Instantiate a converter with overrides (using Postman environment variables)
cfg := swaggman.Configuration{
	PostmanURLBase: "{{RINGCENTRAL_SERVER_URL}}",
	PostmanHeaders: []postman2.Header{
		{
			Key:   "Authorization",
			Value: "Bearer {{my_access_token}}",
		},
	},
}
conv = swaggman.NewConverter(cfg)

// Convert a Swagger spec with a default Postman spec
err := conv.MergeConvert("path/to/swagger.json", "path/to/pman.base.json", "path/to/pman.out.json")
```

## Example

An example conversion is included, [`examples/ringcentral/convert.go`](https://github.com/grokify/swaggman/blob/master/examples/ringcentral/convert.go) which creates a Postman 2.0 spec for the [RingCentral REST API](https://developers.ringcentral.com) using a base Postman 2.0 spec and the RingCentral basic Swagger 2.0 spec.

[A video of importing the resulting Postman collection is available on YouTube](https://youtu.be/5kE4UPXJ-5Q).

Example files include:

* [RingCentral Swagger 2.0 spec](https://github.com/grokify/swaggman/blob/master/examples/ringcentral/ringcentral.spec.swagger2.2019110220191017-1140.json)
* [RingCentral Postman 2.0 base](https://github.com/grokify/swaggman/blob/master/examples/ringcentral/ringcentral.postman2.base.json)
* [RingCentral Postman 2.0 spec](https://github.com/grokify/swaggman/blob/master/examples/ringcentral/ringcentral.spec.postman2.2019110220191017-1140.json) - Import this into Postman

The RingCentral spec uses the following environment variables. The following is the Postman bulk edit format:

```
RC_SERVER_HOSTNAME:platform.devtest.ringcentral.com
RC_APP_KEY:myAppKey
RC_APP_SECRET:myAppSecret
RC_USERNAME:myMainCompanyPhoneNumber
RC_EXTENSION:myExtension
RC_PASSWORD:myPassword
```

For multiple apps or users, simply create a different Postman environment for each.

To set your environment variables, use the Settings Gear icon and then click "Manage Environments"

# Articles and Links

* Medium: [Using Postman, Swagger and the RingCentral API](https://medium.com/ringcentral-developers/using-postman-with-swagger-and-the-ringcentral-api-523712f792a0)
* YouTube: [Getting Started with RingCentral APIs using Postman ](https://youtu.be/5kE4UPXJ-5Q)

 [build-status-svg]: https://github.com/grokify/swaggman/workflows/go%20build/badge.svg
 [build-status-link]: https://github.com/grokify/swaggman/actions
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/swaggman
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/swaggman
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/swaggman
 [docs-godoc-link]: https://pkg.go.dev/github.com/grokify/swaggman
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/swaggman/blob/master/LICENSE.md
