---
title: MCP hosting posture
owner: bifunctional
status: active
updated: 2026-08-07
tags: [mcp, hosting, operations, cost, okf]
---

# MCP hosting posture

## Decision

Use a small, always-on hobby VPS for lightweight organization services when one is needed. Nebius Cloud is the preferred first option. Keep the organization MCP server local-first and stdio-based; a public or remotely reachable MCP endpoint is not part of this decision.

## Intended use

- Run low-traffic supporting services that benefit from being continuously online.
- Use Render for fast demos and temporary deployments.
- Evaluate Railway and Fly.io when their developer experience, regional availability, or managed-service features better fit a particular demo or service.
- Keep persistent data, secrets, private registries, and audit logs out of public repositories and apply the existing access-control policy before any deployment.

## Nebius planning estimates

The following planning estimates are before VAT and should be rechecked in the provider console before provisioning:

| Instance | Disk | Estimated monthly cost, before VAT |
| --- | --- | --- |
| 1 vCPU, 2 GiB RAM | 20 GiB | about $14.85 |
| 1 vCPU, 4 GiB RAM | 20 GiB | about $19.52 |

Start with the smaller instance for a low-traffic service. Increase capacity only after observing memory, CPU, and disk use.

## Budget handling

Cloud-account credits and AI-inference token balances are restricted financial information. Track their exact values in a private operator system, not in this public repository. Before provisioning or enabling inference, confirm the available balance, applicable VAT, expected monthly spend, and a shutdown or spend-limit mechanism.

## Operating guardrails

- Prefer a single small instance and documented deployment configuration over unmanaged production-like infrastructure.
- Do not expose the MCP stdio server through a network listener without a separate, reviewed authenticated-transport design.
- Keep credentials in the provider or deployment secret store, never in repository files or Markdown.
- Record durable provider or architecture changes as an OKF decision document when they materially change the hosting posture.
