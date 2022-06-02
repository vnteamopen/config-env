# config-env

A tool for embedding variable content into the input file.

The tool parse input file, loads variable environment to pattern `{{env "variable_env"}}`, and writes to the output path.

![config-env](https://raw.githubusercontent.com/vnteamopen/config-env/main/config-env.png)

## Features

 - Load input file and replace `{{env "variable_env"}}` with value of variable environemtn `variable_env`

## Installation

### From source

Download the source code and try with:

```
go build -o output/config-env
```

Use `config-env`

### Use from Docker

Pull the docker image from:

```
docker pull ghcr.io/vnteamopen/config-env:main
```

## Quickstart:

1. Create a file `config.yml` with following content:

config.yml
```yml
name: thuc
host: {{env "HOST_NAME"}}
port: {{env "PORT"}}
path: /hello
```

2.1. Run config-env

```bash
HOST_NAME=localhost PORT=1234 ./config-env config.yml output.yml
```

or

```bash
HOST_NAME=localhost PORT=1234 config-env config.yml output.yml
```

2.2. Run config-env with docker

```bash
docker run --rm -it -v $(pwd):/files/ -w /files -e HOST_NAME=localhost -e PORT=1234  ghcr.io/vnteamopen/config-env:main /app/config-env ./config.yml ./output.yml
```

4. output.yml will be write with content

```yml
name: thuc
host: localhost
port: 1234
path: /hello
```

## Future

Check https://github.com/vnteamopen/config-env/issues
