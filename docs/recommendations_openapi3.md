# OpenAPI 3 Spec

This document is a recommended set of minimal OpenAPI 3 spec properties so that a spec can be used by various ecosystem tools such as API References, API Explorers and Client SDK generators.

## Operation

| Property | Requirement | Notes |
|----------|-------------|-------|
| `description` | `MUST` | Used by API References |
| `operationId` | `MUST` | May be used to auto-generate client SDK method names |
| `responses` | `MUST` | Minimally must have `200` response |
| `summary` | `MUST` | Used in API References, such as Swagger UI and ReadMe.io |
| `tags` | `MUST` | Ideally should have 1 and only 1 tag. Tags are used to organize endpoints in API References and Client SDKs. More than 1 tag may not be supported will in some software |