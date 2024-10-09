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

## CLI Plugin Installation

### Manual Installation

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
