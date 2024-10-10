# `docker run`, better

```text
           ⣀⣠⣤⢴⡶⡶⡶⣶⢶⣦⣤⣄⣀
       ⢀⣤⡈⠩⣉⣁⠁       ⠈⠉⠙⠻⢵⣤⡀
     ⣠⢾⠛⠉  ⠈⣛⡟⢷⣦⣄         ⠉⠻⡷⣄
   ⢠⣞⠏⠁    ⠠⣒⠯⣗⢮⣚⢷⣦⡀        ⠈⠻⣷⡀
  ⣰⡯⠁    ⣀⣰⢽⠲⠝⢎⣛⣊⣋⣓⣛⣂⣀        ⠘⢿⠄     ⣀⡀    __________   Docker    ____  ___
 ⣰⡟   ⢀⠤⠌⣃⣭⢶⣞⣿⣽⣾⣻⣽⣯⢿⣽⡯⠿⠶⣶⢦⣤⣤⠤⣤⣤⣤⠴⠶⠒⠚⠉⠉      \______   \__ __  ____ \   \/  /
⢠⡿   ⠐⣩⣶⣻⣽⡾⠿⠙⠓⠉⠈⠁ ⣯   ⠘⢓⣠⣴⣖⣛⣉⡡  ⠰⣾           |       _/  |  \/    \ \     /
⣸⡇  ⣠⣾⣟⠞⠋⠁⣀⣠⣤⣴⣶⡶⡦ ⠻⣄⣀ ⢨⣤⣤⠴⠖⠋⠁    ⡿⡆          |    |   \  |  /   |  \/     \
⣿⡅ ⢰⣟⠇⠁⠤⠒⠉⠉ ⣀⣀⣤⣤⣤⢶⠶⠞⠛⠛⠉⠁         ⡿⡇          |____|_  /____/|___|  /___/\  \
⣷⡇ ⣟⠃  ⢀⣠⡴⠞⠛⠉⠉⣰⡟⣠⠟     ⡀         ⡿⡇                 \/           \/      \_/
⢸⣇ ⢻⡀ ⣰⢻⡁   ⡠⠞⠉⠐⠁   ⣠⡶⠋         ⢀⢿⠃
⠈⣷⡄⠘⢧⡀⠏⠘⢷⣄⣀       ⣠⣾⠏           ⣸⡟
 ⠘⣷⡄⠈⠳⣄⡀ ⠙⠿⣻⣅⡀  ⢀⣼⢿⡞           ⣰⠿
  ⠈⢟⣦ ⢻⠻⠶⢤⣄⣀⠉⠛⠳⠶⣞⣿⣽⠁         ⢀⣴⠟⠁
    ⠙⢷⣄⡀  ⠙⠏⠙⠛⠓⠲⣶⣖⡏        ⢀⡴⡽⠋
      ⠙⠻⣦⣄⡀     ⣷⢿⡅     ⢀⣠⣴⠟⠋
         ⠙⠙⠷⣶⣤⣤⣄⠸⣟⣷⠠⣤⣤⡶⠾⠛⠉
              ⠁⠉ ⠹⣽⡄⠉
                  ⠘⢷⡀
                    ⠁
```

See the [Docker RunX reference](/docs/reference/runx.md) for more information.

## Main Usage

Let's imagine you want to show your container can run a specific action. For instance `echo Hello` from `alpine` image.

Create a `runx.yaml` file with the following content:

```yaml
actions:
  - id: hello # required, the action identifier
    type: run # required, the action type, only `run` is supported for now
    env: # a list of environment variables that needs to be set
      - USER
    cmd: --rm {{.Ref}} echo hello {{env "USER"}} # `.Ref` will be replaced by the reference the user provided
```

And let's create a documentation file called `README.md`:

```markdown
# Hello!

Run the `hello` action with `docker runx` to display a message.
```

Now decorate the `alpine` image and push it under `NAMESPACE/REPOSITORY`:

```
$ docker runx decorate alpine --tag NAMESPACE/REPOSITORY
```

> [!TIP]
> `runx.yaml`  and `README.md` are the default file names.
> You can specify the file names using `--with-config` and `--with-readme` flags.
> It's intended to have both files, if you don't want to provide one or the other, use `--no-config` or `--no-readme` flags.

You can then display the embedded readme:

```
$ dockedr runx NAMESPACE/REPOSITORY --docs
```

Or run the `hello` action:

```
$ docker runx NAMESPACE/REPOSITORY hello
```

See more examples in the [examples](/examples) directory.

## Reference

### Actions Definition

All the possible actions are defined in a yaml file. The file is named `runx.yaml` by default, but you can specify another name using the `--with-config` flag when using the `decorate` command.

```yaml
# Optional.
# Define a default action to run if non is provided by the user.
default: ACTION_ID
# List of all the possible actions for this image.
actions:
    # The action identifier the user can use to run it.
  - id: ACTION_ID
    # Optional.
    # A description of the action, that will be displayed in the documentation or listing.
    desc: DESCRIPTION
    # The type of action to run.
    # If set to `run`, the command will be a `docker run`.
    # If set to `build`, the command will be a `docker buildx build`.
    type: run|build
    # Optional.
    # A list of environment variables that needs to be set before running the command.
    env:
      - ENV_VAR
    # Optional.
    # A list of shell script commands. Their output can be used in the `cmd` field.
    shell:
      NAME: COMMAND
    # Optional.
    # Path to a dockerfile to use with an action of type `build`
    Dockerfile: DOCKERFILE_PATH
    # Optional.
    # A list of options that can be provided by the user.
    opts:
      - name: OPTION_NAME # Name of the option. Also used in the local override or with `--opt` flag.
        desc: DESCRIPTION # Description, rendered in the documentation of the action.
        prompt: PROMPT # A specific prompt to ask the user for the value.
        required: true|false # If required, an empty value will not be accepted.
        values: # A list of possible values for the option. If set, a select will be displayed to the user.
          - VALUE
    # The command to run. It's defined as a Go template and can use the following variables:
    # - `.Ref` will be replaced by the reference to the image the user provided.
    # - `.IsTTY` indicates if the command is run in a TTY environment.
    #   Useful to add the `-t` flag to the `docker run` command: `{{if .IsTTY}}-t{{end}}`
    # - `.Dockerfile` will be replaced by the path to the Dockerfile if one has been provided.
    #
    # In addition, the command can use the following functions:
    # - `{{env "ENV_VAR"}}` will be replaced by the value of the environment variable `ENV_VAR`.
    #   The environment variable needs to be defined in the `env` section.
    # - `{{opt "OPTION"}}` will be replaced by the value of the option `OPTION`.
    #   The value needs to be provided by the local configuration, on the command line or interactively.
    # - `{{sh "COMMAND"}}` will be replaced by the output of the shell command `COMMAND`.
    #   The command will be run using https://github.com/mvdan/sh without a standard input.
    cmd: COMMAND
```

### Local Override

A local file `.docker/runx.yaml` can be used to override the actions defined in the image manifest.
This is useful to configure some actions for a specific project for instance.

```yaml
# Optional.
# It allows to define a default reference to an image if none is provided by the user.
# with the ref set to IMAGE a `docker runx` is equivalent to `docker runx IMAGE`
ref: IMAGE
# Optional.
# Defines some override for one or more images
images:
  # Reference of the image
  IMAGE:
    # Optional.
    # Define a default action to run if none is provided by the user.
    # This overrides the `default` action in the image `runx.yaml` configuration.
    default: ACTION_ID
    # Optional.
    actions:
      # Specify the action to override.
      ACTION_ID:
        opts:
          # Override the value of an option.
          # If set, the option will not be prompt to the user, except if `--ask` flag is used.
          # If `--opt OPTION_NAME=VALUE` is used, this will override the value from this file.
          OPTION_NAME: OPTION_VALUE
```

## Implementation Details

The main idea behind `docker runx` is to attach a specific image manifest to an existing image or image index. This manifest contains the `runx.yaml` and `README.md` files as layers.

When running `docker runx`, the plugin will look for the image manifest and extract the files. It will then execute the action specified by the user or display the documentation.

The image manifest sets the platform as _unkown_ and add a specific annotation to indicate this is a `runx` manifest.

Here is an example of a `runx` based image index:

```json
{
  "schemaVersion": 2, 
  "mediaType": "application/vnd.oci.image.index.v1+json", 
  "manifests": [
    // other manifests...
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json", 
      "size": 533, 
      "digest": "sha256:.....", 
      "platform": {
        "architecture": "unknown", 
        "os": "unknown"
      }, 
      "annotations": {
        "vnd.docker.reference.type": "runx-manifest"
      }
    }
  ]
}
```

The `runx` image manifest will reference at max two layers, one for the `runx.yaml` file and one for the `README.md` file.

```json
{
  "schemaVersion": 2, 
  "mediaType": "application/vnd.oci.image.manifest.v1+json", 
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json", 
    "size": 356, 
    "digest": "sha256:..."
  }, 
  "layers": [
    {
      "mediaType": "application/vnd.runx.config+yaml", 
      "size": 3062, 
      "digest": "sha256:..."
    }, 
    {
      "mediaType": "application/vnd.runx.readme+txt", 
      "size": 946, 
      "digest": "sha256:..."
    }
  ]
}
```

## CLI Plugin Installation

### Manual Installation

<details>
<summary>Instructions</summary>
To install it manually:

- Download the `docker-runx` binary corresponding to your platform from the [latest](https://github.com/eunomie/docker-runx/releases/latest) or [other](https://github.com/eunomie/docker-runx/releases) releases.
- Rename it as
  - `docker-runx` on _Linux_ and _macOS_
  - `docker-runx.exe` on _Windows_
- Copy the binary to the `runx` directory (you might need to create it)
  - `$HOME/.docker/runx` on _Linux_ and _macOS_
  - `%USERPROFILE%\.docker\runx` on _Windows_
- Make it executable on _Linux_ and _macOS_
  - `chmod +x $HOME/.docker/runx/docker-runx`
- Authorize the binary to be executable on _macOS_
  - `xattr -d com.apple.quarantine $HOME/.docker/runx/docker-runx`
- Add the `runx` directory to your `.docker/config.json` as a plugin directory
  - `$HOME/.docker/config.json` on _Linux_ and _macOS_
  - `%USERPROFILE%\.docker\config.json` on _Windows_
  - Add the `cliPluginsExtraDirs` property to the `config.json` file
```
{
	...
	"cliPluginsExtraDirs": [
		"<full path to the .docker/runx folder>"
	],
	...
}
```
</details>

