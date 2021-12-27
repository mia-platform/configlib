<div align="center">

# Config Lib

[![Build Status][github-actions-svg]][github-actions]
[![Go Report Card][go-report-card]][go-report-card-link]
[![GoDoc][godoc-svg]][godoc-link]

</div>

Config lib is a library to handle configuration in your program.
It handles config file and environments variables.

It uses [viper](https://github.com/spf13/viper) and [koanf](https://github.com/knadh/koanf) to handle, respectively, env variables and configuration.

This library handle json in a case sensitive mode.

## Install

```sh
go get -u github.com/mia-platform/configlib
```

## Example usage

### Load file service json configuration - with json schema validation

```go
type Config struct {}

func loadServiceConfiguration(path, fileName string) (Config, error) {
  jsonSchema, err := configlib.ReadFile(configSchemaPath)
  if err != nil {
    return nil, err
  }

  var config ServiceConfig
  if err := configlib.GetConfigFromFile(fileName, path, jsonSchema, &config); err != nil {
    return nil, err
  }

  return config, err
}

// Load service configuration
config, err := loadServiceConfiguration("my/path", "file")
if err != nil {
  log.Fatal(err.Error())
}
```

### Load file service json configuration - without json schema validation

```go
type Config struct {}

func loadServiceConfiguration(path, fileName string) (Config, error) {
  var config ServiceConfig
  if err := configlib.GetConfigFromFile(fileName, path, nil, &config); err != nil {
    return nil, err
  }

  return config, err
}

// Load service configuration
config, err := loadServiceConfiguration("my/path", "file")
if err != nil {
  log.Fatal(err.Error())
}
```

### Get env variables

This feature is deprecated. Please use another lib, like [this](https://github.com/caarlos0/env).

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct,
and the process for submitting pull requests to us.

## Versioning

We use [SemVer][semver] for versioning. For the versions available,
see the [tags on this repository](https://github.com/mia-platform/configlib/tags).

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE)
file for details

[github-actions]: https://github.com/mia-platform/configlib/actions
[github-actions-svg]: https://github.com/mia-platform/configlib/workflows/Test%20and%20build/badge.svg
[godoc-svg]: https://godoc.org/github.com/mia-platform/configlib?status.svg
[godoc-link]: https://godoc.org/github.com/mia-platform/configlib
[go-report-card]: https://goreportcard.com/badge/github.com/mia-platform/configlib
[go-report-card-link]: https://goreportcard.com/report/github.com/mia-platform/configlib
[semver]: https://semver.org/
