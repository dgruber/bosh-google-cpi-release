# Change Log
All releases of the BOSH CPI for Google Cloud Platform will be documented in
this file. This project adheres to [Semantic Versioning](http://semver.org/).

## [25.7.1] - 2017-03-07

### Fixed
- Previously, requests to Google APIs that returned 5xx response codes were
  retried. This change adds retry support to transport errors (net.Error) that
  are known to be temporary.

## [25.7.0] - 2017-03-03

### Fixed
- Various docs changes

### Changed
- Labels may be specified in cloud properties
- CF docs include TCP router
- Use lateset Docker image in CIA
- Support internal load balancer

## [25.6.2] - 2016-12-19

### Fixed
- Improve tests
- Corrects MTU on garden job

### Changed
- Support `ImageURL` for specifying stemcellA
- Concourse docs incude SSL support
- Docs require setting a password instead of using default

## [25.6.1] - 2016-11-02

### Fixed
- Handles large disk sizes

### Changed
- CF docs run under free-tier quota
- Integration tests can use local stemcell

## [25.6.0] - 2016-10-27

### Changed
- Supports `has_disk` method
- Improvements to stemcell pipelines

## [25.5.0] - 2016-10-17

### Changed
- Release uses Go 1.7.1

### Fixed
- A backend service previously could not have multiple instance groups
  with the same name. This release fixes that, and you may have instance
  groups with the same name associated to a backend service.

## [25.4.1] - 2016-09-14

### Fixed
- Tags that are applied by the director on VM create will be truncated to ensure
  they do not violate the 63-char max limit imposed by GCE.

## [25.4.0] - 2016-09-14

### Changed
- When using a custom service account, a default `cloud-platform` scope is used if
  no custom scopes are specified.

## [25.3.0] - 2016-09-02

### Added
- S3 is now a supported blobstore type.

## [25.2.1] - 2016-08-18

### Fixed
- Underscores are replaced with hyphens in metadata that is applied as labels
  to a VM.

### Added
- Complete Concourse installation instructions, including cloud config and Terraform.

## [25.2.0] - 2016-08-18

### Changed
- Any metadata provided by bosh in the `set_vm_metadata` action will also be 
  propagated to the VM as [labels](https://cloud.google.com/compute/docs/label-or-tag-resources),
  allowing sorting and filter in the web console based on job, deployment, etc.

## [25.1.0] - 2016-08-18

### Added
- The `service_account` cloud-config property may now use the e-mail address
  of a custom service account.

## [25.0.0] - 2016-07-25

### Changed
- The `default_zone` config property (in the `google` section of a manifest)
  is no longer supported. The `zone` property must be explicitly set in the
  `cloud_properties` section of `resource_pools` (or `azs` for a cloud-config
  director.)

## [24.4.0] - 2016-07-25

### Fixed
- An explicit region is used to locate a subnet, allowing subnets with the same
  name to be differentiated.

## [24.3.0] - 2016-07-25

### Added
- A `resource_pool`'s manifest can now specify `ephemeral_external_ip` and
  `ip_forwarding` properties, overriding those same properties in a
  manifest's `networks` section.

## [24.2.0] - 2016-07-25

### Added
- This changelog

### Changed
- 3262.4 stemcell

### Fixed
- All tests now use light stemcells

## [24.1.0] - 2016-07-25

### Changed
- Instance tags can be specified in any `cloud_properties` section of a BOSH manifest

### Removed
- The dummy BOSH release is no longer part of the CI pipeline

### Fixed
- Integration tests will use the CI pipeline stemcell rather than requiring an existing stemcell in a project

[25.7.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.6.2...v25.7.0
[25.6.2]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.6.1...v25.6.2
[25.6.1]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.6.0...v25.6.1
[25.6.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.5.0...v25.6.0
[25.5.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.4.1...v25.5.0
[25.4.1]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.4.0...v25.4.1
[25.4.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.3.0...v25.4.0
[25.3.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.2.1...v25.3.0
[25.2.1]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.2.0...v25.2.1
[25.2.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.1.0...v25.2.0
[25.1.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v25.0.0...v25.1.0
[25.0.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.4.0...v25.0.0
[24.4.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.3.0...v24.4.0
[24.3.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.2.0...v24.3.0
[24.2.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24.1.0...v24.2.0
[24.1.0]: https://github.com/cloudfoundry-incubator/bosh-google-cpi-release/compare/v24...v24.1.0
