# Security policy

## Supported surface

The stdio server and the registry schema are the supported deployment surface. A remotely reachable HTTP transport is intentionally not enabled.

## Reporting

Report vulnerabilities privately to the Bifunctional repository maintainers. Do not open a public issue containing exploit details, credentials, private paths, or client data.

## Operator responsibilities

- Protect the local registry and audit log with operating-system permissions.
- Use least-privilege GitHub authentication.
- Keep private repositories in `repos.local.json`.
- Review every requested write and provide confirmation only when its exact target is understood.
- Treat generated documents as drafts until reviewed.

## Explicit non-goals

- Secret storage or secret distribution.
- Unattended destructive repository administration.
- Arbitrary shell execution.
- Blind indexing of source trees or home directories.
- Mirroring private content into this public repository.
