# Release Monitor
Release Monitor is a small program for monitoring the release of new software versions by regularly checking configured online sources, such as GitHub repositories or folders with software packages. The information about the latest version of the software is written as a [Prometheus](https://prometheus.io/) metrics file, which for example can be ingested by the [Node Exporter](https://github.com/prometheus/node_exporter). The following backends for fetching release information is currently supported.

* GitHub repository
* Folder with list of packages
* Free text search

An example configuration file would look like this:
```yaml
update_interval: 43200
metrics_path: /var/lib/prometheus
github:
  - owner: prometheus
    repo: node_exporter
    regexp: v([0-9.]+)
folder:
  - name: ceph
    info: https://ceph.io/community/blog
    path: http://download.ceph.com/tarballs
    regexp: ceph_([0-9.]+).orig.tar.gz
regexp:
  - name: wireguard
    info: https://www.wireguard.com
    path: https://git.zx2c4.com/wireguard-linux-compat/plain/src/version.h
    regexp: (\d+\.\d+\.\d+)
    date: \d+\.\d+.(\d+)
    format: 20060102
```

The `update_interval` in seconds specify how often the program will check the configured sources and the `metrics_path` is where the metrics information file is written. The remaining sections are optional for each different backend, as described below.

## Download
The monitor is written in [Go](https://golang.org) and it can be downloaded and compiled using:
```bash
go get github.com/mrlhansen/release_monitor
```

## Backends
This section contains a short description of the different backends and how they are used in the configuration file.

### GitHub
This backend uses the GitHub API to retrieve information about tags for given repository.
* **owner**: Username of the repository owner.
* **repo**: Name of the repository
* **regexp**: Optional regexp pattern for filtering the version tag. It should contain a single capture group for the version string.

### Folder
This backend scans a list of files/packages using the specified regexp pattern and find the newest version of that package by looking at the associated date.
* **name**: Name of the software.
* **info**: Optional URL for the software, which will be exported as part of the metrics.
* **path**: URL for the list of files.
* **regexp**: Regexp pattern for software package we are checking. It should contain a single capture group for the version string.

### Regexp
This backend does a simple regexp of the entire file and reports the first match as the latest release. There is a regexp both for the version string and the release date, both of which must be found in the file, otherwise no metrics will be generated.
* **name**: Name of the software.
* **info**: Optional URL for the software, which will be exported as part of the metrics.
* **path**: URL for the page/file we are searching in.
* **regexp**: Regexp pattern for the version tag. It should contain a single capture group for the version string.
* **date**: Regexp pattern for the release date. It should contain a single capture group for the date string.
* **format**: Format of the date string in [Go format](https://programming.guide/go/format-parse-string-time-date-example.html).
