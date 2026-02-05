# Security Policy for Repo-lyzer

Repo-lyzer is a developer-focused CLI tool written in Go that analyzes and compares GitHub repositories to help understand code structure, quality, and changes efficiently. This document describes how we handle vulnerability reports and what users and security researchers should expect.

- Current project version: 1.0.6
- Repository: Repo-lyzer
- Language: Go

## Supported Versions

| Version    | Supported |
| ---------- | --------- |
| 1.0.x      | :white_check_mark: (supported) |
| < 1.0      | :x: (not supported) |

We will provide security updates for releases in the 1.0.x series. Older major or pre-1.0 releases are not guaranteed to receive fixes.

## Reporting a Vulnerability

If you believe you have found a security vulnerability in Repo-lyzer, please follow the responsible disclosure process below.

Preferred reporting methods (in order):
1. Email: security@repo-lyzer.dev (PGP: see below) — use this for sensitive or private reports.
2. GitHub Security Advisory: create a confidential security advisory in this repository if you prefer using GitHub’s disclosure workflow.
3. If neither option is available, open a GitHub issue and mark it clearly as a security report; maintainers will respond and, if needed, move it to a private channel.

When emailing, please encrypt sensitive information using our PGP key:
- PGP key / fingerprint: [PROVIDE PGP KEY/FINGERPRINT HERE]

If you cannot encrypt, still send the report but mark carefully and we will respond with a secure channel for further details.

Do not post exploit details publicly until we have had a chance to respond and coordinate disclosure.

## What to include in your report

To help us reproduce and triage the issue quickly, please include:

- Affected Repo-lyzer version(s) (e.g., 1.0.6)
- Platform and environment (OS, Go version, installation method, flags)
- Components affected (CLI, parsing library, comparison logic, etc.)
- Clear, step-by-step reproduction instructions
- Minimal proof-of-concept or sample repository demonstrating the issue (if safe to share)
- Expected vs actual behavior and potential impact (e.g., remote code execution, information disclosure)
- Any logs, stack traces, or error messages
- Your contact details and whether you request anonymity or coordinated disclosure
- If possible, an estimated CVSS score or severity rationale

Use this template (paste into your message or advisory):
- Title:
- Version:
- Component:
- Severity (your assessment):
- Steps to reproduce:
- Proof-of-concept / test repo (if available):
- Mitigations attempted:
- Contact & disclosure preference:

## Our response process & timelines

We aim to follow the timelines below for handling valid reports. These are goals and may vary with complexity and severity.

- Acknowledgement: within 48 hours of receiving a valid report.
- Triage: initial triage within 5 business days (we may request more information).
- Fix / Mitigation:
  - Critical (remote code execution, elevation of privilege, data exfiltration): aim to provide a fix or mitigation within 7 days.
  - High: aim for a fix within 30 days.
  - Medium: aim for a fix or mitigation within 90 days.
  - Low: addressed in a future release; we will communicate expected release windows.
- Updates: if a fix requires longer than indicated, we will provide status updates at least weekly until resolved.
- Disclosure: coordinate public disclosure; we will not publicly disclose without reporter coordination unless required by law.

If a report is received privately and a fix has been released or a coordinated disclosure timeline has elapsed without resolution, we may disclose details to protect users. We will always attempt to give the reporter reasonable notice.

## How we publish security fixes

- Security fixes will be released via GitHub Releases and, where applicable, GitHub Security Advisories.
- Patches will be backported to supported 1.0.x releases when feasible.
- Release notes will contain a high-level description of the issue and a mitigation recommendation. Detailed exploit-proof-of-concept information may be omitted until users are updated or until coordinated disclosure.

## CVE and crediting

- We will coordinate with the reporter about issuing a CVE when appropriate.
- By default, we credit reporters in release notes and acknowledgements unless the reporter requests anonymity.
- If you would like us to withhold credit, state this in your report.

## Safe harbor & legal

- We welcome responsible security research. We do not seek legal action against security researchers who adhere to the reporting guidelines and act in good faith.
- Do not access or modify user data beyond what is necessary to demonstrate the vulnerability. Avoid actions that could cause harm, data destruction, or unacceptable privacy invasion.
- If you are unsure what is reasonable, contact us using the secure channel first.

## Out-of-scope reports

- Bugs that do not affect security (e.g., cosmetic issues or feature requests) should be reported as regular issues.
- Third-party dependency vulnerabilities should be reported to the maintainers of those dependencies; still include them here if you need assistance.

## Disclosure policy example timelines

- Example: Critical RCE reported — Acknowledge within 48 hours; triage in 2 business days; mitigation patch released in <=7 days; coordinated disclosure after patching and user notification.
- Example: Medium information leak — Acknowledge within 48 hours; triage in 5 business days; fix released within 90 days or next minor release with mitigation recommendations provided.

## Contact & escalation

Primary contact:
- agnivamukherjee977@gmail.com
If you do not receive an acknowledgement within 72 hours, follow up to the same address or open a GitHub security advisory and reference your original report.

## Additional notes for maintainers

- Triage checklist: reproduce, classify severity, identify root cause, prepare a fix, prepare tests, prepare release notes, coordinate advisory/CVE, publish patch and notify users.
- Keep a private changelog for security-only changes where appropriate.
- Maintain PGP key and update fingerprint in this document if it changes.

## Policy versioning

- This SECURITY.md is version 1.0 (policy). Last updated: 2026-02-02.
- We may update this policy as the project and team evolve. Always refer to the version in the file header.

Thank you for helping keep Repo-lyzer safe. We appreciate coordinated, responsible disclosure.
