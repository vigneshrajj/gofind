# GoFind

GoFind supercharges your browser address bar by providing short predictable aliases for performing any kind of searches.

Inspired from: [GoLinks](https://git.mills.io/prologic/golinks)

## Features

- Create short aliases for searching your favorite sites or for bookmarking commonly used web pages
- Provides additional aliases for performing common tasks like base64 encoding, decoding, sha256 encoding, etc.
- Stores all your data locally in an sqlite database allowing easy backup for your data
- Lightweight and fast since it is built with technologies like Go and HTMX

## Demo

https://github.com/user-attachments/assets/20d01905-e114-48c3-9845-52deb55af0ee

## Example usages
- Invoke ChatGPT from the address bar by typing 'c' followed by your query: `c how to build a spaceship`. [Command: `#a c https://chatgpt.com/?q=%s`]
- Directly open a specific email inbox by typing `gm personal` or `gm work` or just `gm` which would default to open the personal inbox [Command: `#a gm https://mail.google.com/mail/u/{work:0,personal$(default):1}/#inbox`]
- If you have a long command, you can trigger it by typing part of the command as long as there are no other commands with conflicting name. Eg: Command added this way: `#a commandWithLongAlias https://google.com` can be triggered by typing `com` and pressing enter
- You can run custom bash scripts from the address bar
- You can open any file that can be viewed in a browser like pdf, txt, etc. by opening the file in the browser and then prefixing it with `#a <alias>`, like `#a f file://home/path/to/file.pdf`, then you can use the alias to directly open the file in the browser
- Run multiple commands at once in multiple tabs by separating the commands with `;;`. Example, `g search something in google;;#a alias https://test.com add an alias;;gm work` would run all the commands in 3 different tabs
- And much more:
    - `#a tiny https://tinyurl.com/api-create.php?url={1}` - Use URL shortener service by visiting a website and prefixing the URL with `tiny https://...` to generate a shortened URL
    - `#a insta https://www.instapaper.com/text?u={1}` - View current page in instapaper reader by prefixing the url with `insta https://...`
    - `#a pkt https://getpocket.com/save?url={1}` - Save current page to Pocket by prefixing the url with `pkt https://...`

## Requirements

- Docker or Docker Compose

## Getting Started

### Create a Compose File

#### Docker Compose

Copy the following into a new `docker-compose.yml` file and make any moditications as necessary:

```yml
services:
    gofind:
        image: vigneshrajj/gofind:latest
        ports:
            - 3005:3005
        volumes:
            - ./db:/db
            - ./path/to/local/files:/files # Optional, files located inside this folder can be opened directly using a command
            - ./path/to/bash/scripts:/user_scripts # Optional, executable scripts located inside the below folder can be run with commands
        environment:
            - ENABLE_ADDITIONAL_COMMANDS=false
            - IT_TOOLS_URL=http://localhost:8081 # Optionally, add this url variable along with the below image for enabling IT Tools restart: unless-stopped
    # Optionally, enable my custom IT Tools fork for accessing lots of developer tools right from the address bar
    it-tools:
        image: 'vigneshrajj/it-tools:latest'
        ports:
            - '8081:80'
        restart: unless-stopped
        container_name: it-tools
```

You can run the application using the following command:
```bash
docker compose up -d
```

#### Docker

Alternatively, you can run the following command in your terminal to achieve the same result:

```bash
docker run -d \
  --name gofind \
  -p 3005:3005 \
  -v ./db:/app/db \
  -e ENABLE_ADDITIONAL_COMMANDS=true \
  --restart unless-stopped \
  vigneshrajj/gofind:latest
```

### Set as default search engine

##### Chrome

1. Go to settings by typing `chrome://settings` in the address bar
2. Go to `Search Engine` tab from the left sidebar
3. Click on `Add` button next to **Site Search** and fill the form with these values:
    - Name: GoFind
    - Shortcut: go
    - URL: http://localhost:3005/search?query=%s
4. Click on the three dots icon next to GoFind and choose **Make default**
5. Type `#l` in the address bar to verify if the process was successful

### Additional Configurations

`ENABLE_ADDITIONAL_COMMANDS` - Adds many commonly used search aliases to your database

## Usage

### Commands

- `#l` - Lists all available commands
- `#a <alias> <search_query> <description(optional)>` - Adds a new command
    - Example: `#a g google.com/search?q=%s`
    - Example: `#a g google.com/search?q={1}&q2={2}`
    - Example: `#a gm https://mail.google.com/mail/u/{r:0,vr:1}/#inbox`
    - Example: `#a d https://test.com Some description for the query`
    - Example: `#a file /home/path/to/file.pdf Open a file directly using an alias`
- `#d <alias>` - Deletes an existing command
    - Example: `#d gm`
- `#cmd <scriptname> <arguments>` - Runs a user executable script located at the folder provided in the docker compose file
    - Example: `#cmd bm https://google.com some text as argument` The script named `bm`(or `bm.sh` or with any file extension) would be called with the rest of the string as a single argument. So `https://google.com some text as argument` would be passed entirely as a single string to the `bm` script
- `<alias> <argument1> <argument2> ...` - Searches the website denoted by the alias along with the provided arguments
    - Example: `g how to build a spaceship`
    - Example: `gm vr`

##### Types of Arguments

- **Search String arguments** -  `#a g google.com/search?q=%s`
    - Use alias g followed by any number of arguments for searching
    - `g how to make a hello world program in go`
- **Numbered arguments** -  `#a g google.com/search?q={1}&q2={2}`
    - Use alias g followed by 2 arguments and each of them will be placed at the respective places for searching
    - `g abc efg` would become `https://google.com/search?q=abc&q2=efg`
- **Key Value arguments** - `#a gm https://mail.google.com/mail/u/{r:0,vr:1}/#inbox`
    - Use alias gm followed by specific key from the keys provided in the command (keys from above example are r, vr) and it will be replaced by the corresponding value (values from above example are 0, 1):
        - `gm vr` would become `https://mail.google.com/mail/u/1/#inbox`
    - Mark a value as default value if no arguments are passed:
        - `#a gm https://mail.google.com/mail/u/{r$(default):0,vr:1}/#inbox` - Creating a command with `$(default)` before the colon specifies it as the default
        - For the above command, `gm` would become `https://mail.google.com/mail/u/0/#inbox`

*Note:* You cannot use multiple types of arguments in a single command. So `#a g https://google.com/%s?q={1}&q2={key:val}` may not work.

##### Additional Utilities

Some examples for additional utilities provided by IT Tools include:

| Utility          | Alias  | Example                      |
|------------------|--------|------------------------------|
| Hash Text        | hash   | `!hash abcd`                 |
| Base64 Encoding  | b64    | `!b64 true abcd`             |
| Base64 Decoding  | d64    | `!b64 false <base64 string>` |

## Development

### Running Locally

Built with Go v1.23.2
- Clone the repository:
```bash
git clone https://github.com/vigneshrajj/gofind && cd gofind
```
- Install dependencies:
```bash
go mod download
```
- Run the application
```bash
go run ./cmd/gofind/main.go
```
The application will start on http://localhost:3005.

### Running Tests

- Run the following command to run the tests:
```bash
go test ./tests
```

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or features.

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Author

Vignesh Raj
