# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.0]

Helm charts for every service, the Silver to OpenGovMail rename, ThunderID 1.0.1 with support for an external ThunderID, Nextcloud with ThunderID single sign-on, and versioned images built on release tags. This release moves from `v.X.Y.Z` tags to `vX.Y.Z`.

## What's Changed
* Update project status of features in README by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/308
* Add Thunder API Consumer Contract Guide for Silver by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/316
* 309 task add email application setup step in thunder by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/312
* Add example Thunderbird extension for OAuth2 authentication by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/311
* Create Helm subchart for OpenDKIM service with configuration and templates by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/320
* Add Helm subchart for cert-manager to manage TLS certificates by @maneesha-xyz in https://github.com/OpenGovMail/OpenGovMail/pull/325
* Upgrade Thunder to version 0.32.0 and update configuration for console application by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/328
* Disable SMTP VRFY and EXPN commands for security by @manupawickramasinghe in https://github.com/OpenGovMail/OpenGovMail/pull/332
* #333 Update gen-postfix-conf.sh and replace hardcoded Thunder credentials  by @manupawickramasinghe in https://github.com/OpenGovMail/OpenGovMail/pull/334
* #63 Standardize Docker Compose service names by @prdai in https://github.com/OpenGovMail/OpenGovMail/pull/337
* docs: remove unused conf/users.yaml and update Adding users to Thunder flow (#338) by @prdai in https://github.com/OpenGovMail/OpenGovMail/pull/339
* 322 task create helm subchart for rspamd service by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/324
* Feat/postfix charts by @maneesha-xyz in https://github.com/OpenGovMail/OpenGovMail/pull/342
* feat(charts): add Thunder identity server via umbrella dependency by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/344
* Feat/raven seaweed charts (#343) by @maneesha-xyz in https://github.com/OpenGovMail/OpenGovMail/pull/345
* chore: docs by @maneesha-xyz in https://github.com/OpenGovMail/OpenGovMail/pull/346
* Wip/silver charts by @maneesha-xyz in https://github.com/OpenGovMail/OpenGovMail/pull/347
* Rename Silver → OpenGovMail across the codebase by @saknarajapakshe in https://github.com/OpenGovMail/OpenGovMail/pull/352
* fix: update references from Silver to OpenGovMail in configuration files by @saknarajapakshe in https://github.com/OpenGovMail/OpenGovMail/pull/360
* fix: fail fast on missing/unparseable domain config in OpenDKIM and Postfix entrypoints by @Dushmilan in https://github.com/OpenGovMail/OpenGovMail/pull/361
* docs(charts): fix Thunder admin credentials instructions in README by @saknarajapakshe in https://github.com/OpenGovMail/OpenGovMail/pull/359
* feat(charts): migrate to thunderid 1.0.1 with a raven service account by @maneesha-xyz in https://github.com/OpenGovMail/OpenGovMail/pull/358
* metadata-service: validate PUSH_INTERVAL_SECONDS and parse ENABLE_PUS… by @ruchitasingla in https://github.com/OpenGovMail/OpenGovMail/pull/377
* chore(docker): update and pin service base images by @Dushmilan in https://github.com/OpenGovMail/OpenGovMail/pull/378
* feat(charts): add Nextcloud with ThunderID single sign-on by @maneesha-xyz in https://github.com/OpenGovMail/OpenGovMail/pull/372
* feat(charts): support an external ThunderID via global.thunder.url by @maneesha-xyz in https://github.com/OpenGovMail/OpenGovMail/pull/381
* feat: build versioned images on release tags and pin chart images by @maneesha-xyz in https://github.com/OpenGovMail/OpenGovMail/pull/382


**Full Changelog**: https://github.com/OpenGovMail/OpenGovMail/compare/v.0.3.0...v0.4.0

## [0.3.0]

Multi-domain support, improved configuration handling, group address validation, Thunder upgrade to 0.27.0, and removal of legacy socketmap/Thunder integration.

## What's Changed
* feat: enhance setup scripts for multi-domain support and improve configuration handling by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/301
* feat: add support for group address validation in user-exists map and implement group validation logic by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/302
* feat: upgrade Thunder image version to 0.27.0 in Docker Compose by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/303
* Refactor: Remove socketmap handler and thunder integration by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/305


**Full Changelog**: https://github.com/OpenGovMail/OpenGovMail/compare/v.0.2.1...v.0.3.0

## [0.2.1]

This release focuses on making improvements to the observability of OpenGovMail.

## What's Changed
### Added
* Update thunder idp to latest version to support admin initiated user registration via url flow by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/265
* Move observability services into mail stack and update configur… by @maneeshaxyz in https://github.com/OpenGovMail/OpenGovMail/pull/264
* Add Grafana domain update functionality in gen-observability script by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/275

### Removed
* Remove Change Password UI by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/267

### Fixed
* Fix observability bugs by @maneeshaxyz in https://github.com/OpenGovMail/OpenGovMail/pull/270
* Fixed the issue of not saving the attachments in the blob storage. by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/268
* Fix observability port issue by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/273
* Loki not working correctly issue resolved. by @maneeshaxyz in https://github.com/OpenGovMail/OpenGovMail/pull/277
* Secure the rspamd web UI dashboard by @Aravinda-HWK in https://github.com/OpenGovMail/OpenGovMail/pull/279

**Full Changelog**: https://github.com/OpenGovMail/OpenGovMail/compare/v.0.2.0...v.0.2.1

## [0.2.0] - Raven Integration and Platform Hardening

This release introduces Raven-based mail architecture, major security and
observability improvements, and stronger operational tooling.

### Added

- Raven delivery integration and related configuration updates
  (`#120`, `#121`, `#123`).
- Change password Web UI and backend services (`#201`, `#203`, `#206`, `#207`).
- Load testing suite, utilities, and CI jobs (`#211`, `#213`, `#215`, `#229`,
  `#230`, `#232`, `#234`).
- Observability stack with Prometheus, Loki, Promtail, and Grafana dashboards
  (`#173`, `#174`, `#178`, `#181`, `#183`, `#184`, `#186`, `#188`, `#190`,
  `#198`).
- Smart attachment/blob storage support (`#238`).
- Docker cleanup and service control scripting improvements (`#121`, `#133`).

### Changed

- Replaced Dovecot-based authentication and retrieval flows with Raven services
  and database-backed transport (`#138`, `#141`, `#152`).
- Migrated service configuration layout and script paths toward `conf/` and YAML
  driven generation (`#129`, `#137`).
- Improved container build/push workflows and integrated image scanning
  refinements (`#131`, `#144`).
- Enhanced multi-domain certificate/domain handling and DKIM script workflows
  (`#157`, `#158`).
- Removed Unix socket usage from core OpenGovMail services (`#236`).

### Fixed

- Resolved ClamAV daily refresh OOM issues (`#143`).
- Removed OCSP stapling configuration that caused compatibility issues (`#168`).
- Corrected detect-changes workflow repository filtering (`#217`).
- Fixed multiple script/config path issues in setup and documentation (`#129`,
  `#137`, `#246`).
- Fixed warning output in Docker containers (`#242`).

### Security

- Increased TLS hardening and added automated TLS security tests (`#166`).
- Removed unnecessary exposed ports for internal services, including SeaweedFS
  hardening (`#141`, `#243`).
- Added centralized ClamAV signature distribution and service cleanup (`#241`,
  `#245`).
- Removed encrypted password exposure in rspamd-related flow (`#176`).

### Documentation

- Expanded and improved README content, badges, and release docs (`#159`, `#225`,
  `#226`, `#234`, `#247`, `#248`).
- Refined contribution and repository guidance docs (`#160`, `#161`).

### Full Changelog

https://github.com/OpenGovMail/OpenGovMail/compare/v0.1.0...v0.2.0

## [0.1.0] - First Stable OpenGovMail M1 Release

This release saves a stable working version of OpenGovMail M1 before integrating
Raven into the repository.

### Added

- Initial Dockerized OpenGovMail M1 services and SMTP host-to-container delivery
  testing (`#16`, `#17`, `#27`).
- Initial API/Web UI setup and operation scripts for running the mail server
  (`#18`, `#31`, `#32`).
- Bootstrap/init automation and user-facing tooling, including user listing and
  quota controls (`#28`, `#29`, `#44`, `#48`).
- Load testing support and related tooling (`#42`).
- Dynamic virtual user provisioning with POP3 support (`#61`, `#65`, `#66`).
- Repository collaboration templates for issues and pull requests (`#68`,
  `#69`).

### Changed

- Consolidated container layout, improved environment defaults, and completed
  service cleanup/hardening (`#21`, `#22`, `#64`, `#75`, `#83`).
- Migrated service configuration toward YAML-driven flows and parsing (`#38`,
  `#39`, `#40`).
- Refactored Thunder endpoints and restricted external access to Thunder IDP
  (`#45`, `#47`, `#58`).
- Updated runtime/logging behavior for Docker compatibility and merged long-lived
  branches to stabilize the release line (`#57`, `#81`, `#85`).
- Expanded and refined project documentation (`#19`, `#30`, `#34`, `#35`,
  `#36`, `#37`, `#55`, `#56`).

### Fixed

- Fixed user initialization and add-user script behavior (`#41`, `#43`).
- Fixed LMTP mail receiving flow and other release-blocking defects (`#78`,
  `#79`).

### Full Changelog

https://github.com/OpenGovMail/OpenGovMail/releases/tag/v0.1.0
