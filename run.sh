#!/bin/bash
[ -d bin ] || mkdir bin &&
go build -o bin/bookings ./cmd/web/*.go &&
./bin/bookings