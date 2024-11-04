# GoFind

GoFind supercharges your browser address bar by providing short predictable aliases for performing any kind of searches. 

## Features

- Create short aliases for searching your favorite sites or for bookmarking commonly used web pages
- Provides additional aliases for performing common tasks like base64 encoding, decoding, sha256 encoding, etc.
- Stores all your data locally in an sqlite database allowing easy backup for your data
- Lightweight and fast since it is built with technologies like Go and HTMX

## Requirements
- Docker Compose

## Getting Started

### Create a Compose File

Copy the following into a new `docker-compose.yml` file and make any moditications as necessary:

```yml
services:
  gofind:
    image: vigneshrajj/gofind:latest
    ports:
      - 3005:3005
    volumes:
      - ./db:/app/db
    environment:
      - ENABLE_ADDITIONAL_COMMANDS=true
    restart: unless-stopped
```

### Run the Application

You can run the application using the following command:
```bash
docker compose up -d
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
- `#a <alias> <search_query>` - Adds a new command
    - Example: `#a g google.com/search?q=%s`
    - Example: `#a g google.com/search?q={1}&q2={2}`
    - Example: `#a gm https://mail.google.com/mail/u/{r:0,vr:1}/#inbox`
- `#d <alias>` - Deletes an existing command
    - Example: `#d gm`
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
    - Use alias gm followed by specific key from the keys provided in the command (keys from above example are r, vr) and it will be replaced by the corresponding value (values from above example are 0, 1)
    - `gm vr` would become `https://mail.google.com/mail/u/1/#inbox`

##### Additional Utilities

| Utility          | Alias  | Example               |
|------------------|--------|-----------------------|
| SHA 256 Encoding | sha256 | `sha256 abcd`         |
| Base64 Encoding  | b64    | `b64 abcd`            |
| Base64 Decoding  | d64    | `d64 <base64 string>` |

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
