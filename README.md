# Maven Resource

[![Go Report Card](https://goreportcard.com/badge/github.com/garethjevans/maven-resource)](https://goreportcard.com/report/github.com/garethjevans/maven-resource)

A concourse resource that can track information about a maven dependency.

## Source Configuration

* `repository`: *Required.* The location of the repository.

* `groupId`: *Required.* The groupId of the artifact to download.

* `artifactId`: *Required.* The artifactId of the artifact to download.

* `type`: *Optional.* The type of the artifact to download.

* `username`: *Required.* Username for accessing an authenticated repository.

* `password`: *Optional.* Password for accessing an authenticated repository.


## Behavior

### `check`: Check for new versions of the artifact.

Checks for new versions of the artifact by retrieving the `maven-metadata.xml` from
the repository.

NOTE: Only valid semantic versions are supported.

### `in`: Fetch an artifact from a repository.

Download the artifact from the repository.

#### Additional files populated

* `version`: The version of the downloaded artifact.

* `filename`: The filename of the downloaded file.

* `uri`: The full uri that can be used to reference the current version of the artifact.

* `sha1`: The sha1 sum of the downloaded file.

* `sha256`: The sha256 sum of the downloaded file. If this is not available from the maven 
   repository, it is calculated from the downloaded file after the `sha1` is validated.

* `cpe`: The version to be used in a CPE referrence.

* `purl`: The version to be used in a PURL reference.
 
### `out`:

Not Implemented.

## Examples

Resource configuration for an authenticated repository:

``` yaml
resource_types:
- name: maven-resource
  type: registry-image
  source:
    repository: ghcr.io/garethjevans/maven-resource
    tag: latest

resources:
- name: artifact
  type: maven-resource
  source:
    repository: https://myrepo.example.com/repository/maven-releases
    artifactId: example-webapp
    groupId: com.example
    type: jar
    username: myuser
    password: mypass
```


Retrieve an artifact and push to Cloud Foundry using [cf-resource](https://github.com/concourse/cf-resource)

``` yaml
jobs:
- name: deploy
  plan:
  - get: source-code
  - get: artifact
    trigger: true
  - put: cf
    params:
      manifest: source-code/manifest.yml
      path: artifact/example-webapp-*.jar
```
