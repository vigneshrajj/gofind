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
- `#a <alias> <search_string>` - Adds a new command
    - Example: `#a g google.com/search?q=%s`
    - Example: `#a g google.com/search?q={1}&q2={2}`
    - Example: `#a gm https://mail.google.com/mail/u/{r:0,vr:1}/#inbox`
- `#d <alias>` - Deletes an existing command
    - Example: `#d gm`
- `<alias> <argument1> <argument2> ...` - Searches the website denoted by the alias along with the provided arguments
    - Example: `g how to build a spaceship`
    - Example: `gm vr`

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
