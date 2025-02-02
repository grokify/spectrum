Spectrum - OpenAPI Spec SDK and Postman Converter
=================================================

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![LOC][loc-svg]][repo-url]
[![License][license-svg]][license-url]

Spectrum is a multi-purpose OpenAPI Spec SDK that includes enhanced Postman conversion. Most of the OpenAPI Spec SDK is designed to support OAS3. Some functionality for OAS2 exists.

The following article provides an overview of OpenAPI spec to Postman conversion:

1. [Blog Introduction](https://medium.com/ringcentral-developers/using-postman-with-swagger-and-the-ringcentral-api-523712f792a0)

## Use Cases

1. Converting from OpenAPI Spec to Postman where `x-tagGroups` are converted to Postman folders.
1. Programmatically editing specs for reproducible results, typically when source spec is maintained.
1. Merging multiple specs when a single spec is desired.
1. Splitting a spec by tags when multiple specs are desired.

## Packages and Major Features

* openapi2 ([godoc](https://pkg.go.dev/github.com/grokify/spectrum/openapi2))
  1. Support for OpenAPI 2 files, including serialization, deserialization, and validation.
  1. Merging of multiple specs
  1. Postman 2 Collection conversion
* openapi3 ([godoc](https://pkg.go.dev/github.com/grokify/spectrum/openapi3))
  1. Support for OpenAPI 3 files, including serialization, deserialization, and validation.
  1. Merging of multiple specs
  1. Splitting specs by tag
  1. Output of spec to tabular format to HTML (API Registry), CSV, XLSX. HTML API Registry has a bonus feature that makes each line clickable. Click any line here: http://ringcentral.github.io/api-registry/
  1. Programmatic API to modify OpenAPI specs using rules
  1. [Programmatic ability to "fix" spec, e.g. change response Content Type to match output (needed for Engage Voice)](docs/openapi3_fix.md)
  1. [OpenAPI 3 linter](openapi3/openapi3lint)
  1. Statistics: Counts operations, schemas, properties & parameters (with and without descriptions), etc.
  1. Postman 2 Collection conversion
  1. Ability to merge in Postman request body examples into Postman 2 Collection
  1. Functionality is built on *kin-openapi*: https://github.com/getkin/kin-openapi
* openapi3edit ([godoc](https://pkg.go.dev/github.com/grokify/spectrum/openapi3edit))
  1. Programmatic SDK-based editor for OAS3 specifications.
* openapi3lint ([godoc](https://pkg.go.dev/github.com/grokify/spectrum/openapi3lint))
  1. Extensible linter for OAS3 specifications.
* postman2 ([godoc](https://pkg.go.dev/github.com/grokify/spectrum/postman2))
  1. Support for Postman 2 Collection files, including serialization and deserialization.
  1. CLI and library to Convert OpenAPI Specs to Postman Collection
  1. Add Postman environment variables to URLs, e.g. Server URLs like `https://{{HOSTNAME}}/restapi`
  1. Add headers, such as environment variable based Authorization headers, such as `Authorization: Bearer {{myAccessToken}}`
  1. Utilize baseline Postman collection to add Postman-specific functionality including Postman `prerequest` scripts.
  1. Add example request bodies, e.g. JSON bodies with example parameter values.
* raml08
  1. Support for parsing RAML v0.8
  1. Limited functionality to extracting OpenAPI v3 `description` and `summary` from `description` and `displayName` respectively.

## Notes

* Postman 4.10.7 does not natively support JSON requests so request bodies need to be entered using the raw body editor. A future task is to add Swagger request examples as default Postman request bodies.
* Postman 2.0 spec supports polymorphism and doesn't have a canonical schema. For example, the `request.url` property can be populated by a URL string or a URL object. Spectrum uses the URL object since it is more flexible. The function `simple.NewCanonicalCollectionFromBytes(bytes)` can be used to read either a simple or object based spec into a canonical object spec.
* This has only been used on the RingCentral Swagger spec to date but will be used for more in the future. Please feel free to use and contribute. Examples are located in the `examples` folder.

## Installation

The following command will install the executable binary `spectrum` into the `~/go/bin` directory.

```bash
$ go get github.com/grokify/spectrum
```

## Usage

### Simple Usage

```
// Instantiate a converter with default configuration
conv := spectrum.NewConverter(spectrum.Configuration{})

// Convert a Swagger spec
err := conv.Convert("path/to/swagger.json", "path/to/pman.out.json")
```

### Usage with Features

The following can be added which are especially useful to use with environment variables.

* Custom Hostname
* Custom Headers

```
// Instantiate a converter with overrides (using Postman environment variables)
cfg := spectrum.Configuration{
	PostmanURLBase: "{{RINGCENTRAL_SERVER_URL}}",
	PostmanHeaders: []postman2.Header{
		{
			Key:   "Authorization",
			Value: "Bearer {{my_access_token}}",
		},
	},
}
conv = spectrum.NewConverter(cfg)

// Convert a Swagger spec with a default Postman spec
err := conv.MergeConvert("path/to/swagger.json", "path/to/pman.base.json", "path/to/pman.out.json")
```

### Example

An example conversion is included, [`examples/ringcentral/convert.go`](https://github.com/grokify/spectrum/blob/master/examples/ringcentral/convert.go) which creates a Postman 2.0 spec for the [RingCentral REST API](https://developers.ringcentral.com) using a base Postman 2.0 spec and the RingCentral basic Swagger 2.0 spec.

[A video of importing the resulting Postman collection is available on YouTube](https://youtu.be/5kE4UPXJ-5Q).

Example files include:

* [RingCentral Swagger 2.0 spec](https://github.com/grokify/spectrum/blob/master/examples/ringcentral/ringcentral.spec.swagger2.2019110220191017-1140.json)
* [RingCentral Postman 2.0 base](https://github.com/grokify/spectrum/blob/master/examples/ringcentral/ringcentral.postman2.base.json)
* [RingCentral Postman 2.0 spec](https://github.com/grokify/spectrum/blob/master/examples/ringcentral/ringcentral.spec.postman2.2019110220191017-1140.json) - Import this into Postman

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

## Articles and Links

* Medium: [Using Postman, Swagger and the RingCentral API](https://medium.com/ringcentral-developers/using-postman-with-swagger-and-the-ringcentral-api-523712f792a0)
* YouTube: [Getting Started with RingCentral APIs using Postman](https://youtu.be/5kE4UPXJ-5Q)

 [build-status-svg]: https://github.com/grokify/spectrum/workflows/test/badge.svg
 [build-status-url]: https://github.com/grokify/spectrum/actions/workflows/test.yaml
 [lint-status-svg]: https://github.com/grokify/spectrum/workflows/lint/badge.svg
 [lint-status-url]: https://github.com/grokify/spectrum/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/spectrum
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/spectrum
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/spectrum
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/spectrum
 [loc-svg]: https://tokei.rs/b1/github/grokify/spectrum
 [repo-url]: https://github.com/grokify/spectrum
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/spectrum/blob/master/LICENSE
