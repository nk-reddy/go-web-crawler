#!/usr/bin/env bash
go build -o crawler && ./crawler "https://crawler-test.com/" 3 25
