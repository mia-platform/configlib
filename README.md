<div align="center">

# Config Lib

[![Build Status][github-actions-svg]](gitub-actions)
[![Go Report Card][go-report-card]](go-report-card-link)
[![GoDoc][godoc-svg]](godoc-link)

</div>

Config lib is a library to handle configuration in your program.
It handles config file and environments variables.

It uses [viper](https://github.com/spf13/viper) and [koanf](https://github.com/knadh/koanf) to handle, respectively, env variables and configuration.

This library handle json in a case sensitive mode.

## Install

```sh
go get -u github.com/mia-platform/config-lib
```

## Example usage

### Get env variables

```go
type EnvVariables struct {
  StringValue string
  BoolValue   bool
}
var envVariablesConfig = []configlib.EnvConfig{
  {
		Key:          "SOME_STRING",
		Variable:     "StringValue",
		DefaultValue: "",
	},
	{
		Key:      "BOOLEAN_KEY",
		Variable: "BoolValue",
		Required: true,
	},
}

var env EnvVariables
err := configlib.GetEnvVariables(envVariablesConfig, &env)
if err != nil {
  panic(err.Error())
}
```

### Load file service json configuration - with json schema validation

```go
type Config struct {}

func loadServiceConfiguration(path, fileName string) (Config, error) {
	jsonSchema, err := configlib.ReadFile(configSchemaPath)
	if err != nil {
		return nil, err
	}
	var config ServiceConfig
	err = configlib.GetConfigFromFile(fileName, path, jsonSchema, &config)
	if err != nil {
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
	jsonSchema, err := configlib.ReadFile(configSchemaPath)
	if err != nil {
		return nil, err
	}
	var config ServiceConfig
	err = configlib.GetConfigFromFile(fileName, path, nil, &config)
	if err != nil {
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

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct,
and the process for submitting pull requests to us.

## Versioning

We use [SemVer][semver] for versioning. For the versions available,
see the [tags on this repository](https://github.com/mia-platform/terraform-google-project/tags).

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE.md](LICENSE.md)
file for details

[github-actions]: https://github.com/mia-platform/config-lib/actions
[github-actions-svg]: https://github.com/mia-platform/config-lib/workflows/Test%20and%20build/badge.svg
[godoc-svg]: https://godoc.org/github.com/mia-platform/config-lib?status.svg
[godoc-link]: https://godoc.org/github.com/mia-platform/config-lib
[go-report-card]: https://goreportcard.com/badge/github.com/mia-platform/config-lib
[go-report-card-link]: https://goreportcard.com/report/github.com/mia-platform/config-lib
