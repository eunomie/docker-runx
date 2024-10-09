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
