# ShortURL app
ShortURL is a simple web application that allows users to create short URLs using a hash function. The shortened URL will redirect to the original URL when accessed.

## Getting Started
To use the ShortURL app, you need to have `Go` installed on your computer. You can download `Go` from the [official website](https://golang.org/dl/).

After installing `Go`, clone the repository using the following command:
```bash
git clone https://github.com/gennadis/shorturl.git
```

Change into the project directory:
```bash
cd shorturl
```

Then, build the app:
```go
go build
```

To start the app, run:
```bash
./shorturl
```

By default, the app will start listening on `localhost` port `8080`.

## Usage
To create a short URL, send a `POST` request to the root endpoint with the URL you want to shorten in the request body. The app will return the shortened URL as plain text in the response body.

To access the original URL, access the shortened URL returned in the previous step. The app will redirect to the original URL.

## Dependencies
The `ShortURL` app has the following dependencies:
- go-chi/chi/v5
- go-chi/chi/v5/middleware

These dependencies can be installed using the go get command:
```go
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware
```

## Configuration
The app can be configured using the config package in internal/`config/config.go`. The default configuration sets the app to listen on port `8080` and use a hash length of `6`. You can modify the configuration by changing the values in `config.go`.

## Contributing
If you would like to contribute to the `ShortURL` app, please create a `pull request` with your changes. All contributions are welcome!
