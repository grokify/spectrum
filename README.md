Swagger2Postman in Go
=====================

`swagger2postman` creates Postman 2.0 specs from a Swagger 2.0 spec.

# Features

* Build on a base Postman 2.0 spec that contains information not readily stored in Swagger Spec to support Postman features such as scripts.
* Supports override parameters such as a Postman environment parameter for the URL hostname, e.g. https://{{MY_HOSTNAME}}/rest
* Supports additional headers, e.g. authorization headers using enviroment variables, e.g. `Authorization: Bearer {{myAccessToken}}`

# Example Usage

There is an example conversion available, `examples/ringcentral/convert.go` which creates a Postman 2.0 spec for the [RingCentral REST API](https://developers.ringcentral.com) using a base Postman 2.0 spec and the RingCentral basic Swagger 2.0 spec.

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
