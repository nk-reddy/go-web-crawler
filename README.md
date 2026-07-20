# Web Crawler

A concurrent Go crawler that follows links on a site and prints the pages it finds.

## Run

Requires Go.

```sh
go mod download
go run . <url> <max-concurrency> <max-pages>
```

Example:

```sh
go run . https://crawler-test.com/ 3 25
```

You can run the included example:

```sh
chmod +x run.sh
./run.sh
```
