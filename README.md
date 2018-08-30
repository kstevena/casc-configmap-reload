# Kubernetes ConfigMap Reload

[![license](https://img.shields.io/github/license/jimmidyson/configmap-reload.svg?maxAge=2592000)](https://github.com/kstevena/casc-configmap-reload)
[![Docker Stars](https://img.shields.io/docker/stars/jimmidyson/configmap-reload.svg?maxAge=2592000)](https://hub.docker.com/r/kstevena/casc-configmap-reload/)
[![Docker Pulls](https://img.shields.io/docker/pulls/jimmidyson/configmap-reload.svg?maxAge=2592000)](https://hub.docker.com/r/kstevena/casc-configmap-reload/)
[![CircleCI](https://img.shields.io/circleci/project/jimmidyson/configmap-reload.svg?maxAge=2592000)](https://circleci.com/gh/kstevena/casc-configmap-reload)

**casc-configmap-reload** is a simple binary to trigger a reload of the Configuration-as-Code Jenkins plugin configuration when Kubernetes ConfigMaps are updated.
It watches mounted volume dirs and notifies the target process that the config map has been changed.
It is available as a Docker image at https://hub.docker.com/r/kstevena/casc-configmap-reload

### Usage

```
Usage of ./out/configmap-reload:
  -jenkins-url string
        the jenkins url
  -password string
        the jenkins password
  -username string
        the jenkins username
  -volume-dir value
        the config map volume directory to watch for updates; may be used multiple times
  -webhook-status-code int
        the HTTP status code indicating successful triggering of reload (default 200)
```

### License

This project is [Apache Licensed](LICENSE.txt)

### Credits

Thanks to [jimmidyson](https://github.com/jimmidyson) for his wonderfull [configmap-reload](https://github.com/jimmidyson/configmap-reload) project
