# Spring Initializr Resource

Tracks the available versions of [Spring Boot](http://spring.io/projects/spring-boot)
on a given [Spring Initializr](https://github.com/spring-io/initializr/) instance

## Source Configuration

All fields are optional and have reasonable defaults where necessary.

* `url`: The base URL of the Initializr site. Must be a full URL. Defaults to
https://start.spring.io

* `skip_tls_validation`: If the URL is pointing to an HTTPS site and this is true,
TLS certificates will not be validated. If using a `url` with self-signed certs,
consider using `ca_certs` instead.

* `ca_certs`: A list of PEM-encoded X509 trusted certificates that will be used to
establish trust during requests

* `product_version`: A regular expression to match product version e.g. `1\.15\..*`.
Empty values match all product versions.

* `include_snapshots`: If true, include builds that end with `BUILD-SNAPSHOT` in
addition to builds that end in `RELEASE`

* `https_proxy`: A Proxy server URL to use for HTTPS requests. Can have a scheme of either
`http`, `https`, or `socks5`

* `http_proxy`: A Proxy server URL to use for HTTP requests. Can have a scheme of either
`http`, `https`, or `socks5`

* `no_proxy`: A comma-separated list of hosts, IPs, and domain names that do not
use one of the above-mentioned proxies.

## Behavior

### `check`: Watch for new versions of Spring Boot available on the initializr

Versions returned will match `product_version`, if set, and will only include versions
that end in `.RELEASE`, unless `include_snapshots` is truthy.

### `in`: Generate a project from the initializr

Will generate a new Spring Boot project (or build file) with the given parameters.

It will place the following files in the bucket:
* `(filename)`: The generated file. If `type` is `maven-project` or `gradle-project`,
  the file name will be `starter.zip`. If `type` is `maven-build`, the file will be
  `pom.xml`, and if it's `gradle-build`, the file will be `build.gradle`.

* `version`: The version of Spring Boot used to generate the project

* `url`: The URL used to generate the project

* `available_dependencies`: A list of all libraries that can be used with this version
  of Spring Boot.

#### Parameters

All fields are optional and have reasonable defaults where necessary.

* `type`: The type of file to generate. Valid options are `maven-project` (default),
  `gradle-project`, `maven-build`, or `gradle-build`

* `dependencies`: A comma-separated list of dependencies to be included in the project

* `packaging`: `jar` or `war` (defaults to `jar`)

* `jdk_version`: What version of Java to use for the project. Currently supported versions
  are `1.8` (default) and `10`

* `language`: The language to generate code in. Current options are `java` (default), `groovy`,
  and `kotlin`

* `group_id`: The Maven group ID to use. Defaults to `com.example`

* `artifact_id`: The Maven artifact ID to use. Defaults to `demo`.

* `version`: The maven version to use (defaults to `0.0.1-SNAPSHOT`)

* `name`: The name of the project (defaults to `demo`).

* `description`: The longer description of the project for the build file.

* `package_name`: The Java package name for the generated code. Defaults to `com.example`

### `out`

Because the initializr is a read-only API, so too is this resource type. If that changes,
functionality will be added.

## Example Configuration

### Resource
```yaml
resource_types:
- name: spring-initializr
  type: docker-image
  source:
    repository: jghiloni/spring-initializr-resource

- name: start-spring-io
  type: spring-initializr
  source:
    url: https://start.spring.io
    product_version: 1\.5\..*
    include_snapshots: false
```

### Plan
```yaml
- get: start-spring-io
  params:
    type: maven-project
    dependencies: data,web,security,actuator
    group_id: com.myco.myproject
    artifact_id: some-cool-app
    package_name: com.myco.myproject.myapp
```
