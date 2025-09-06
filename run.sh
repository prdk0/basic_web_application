#!/bin/bash

go build -o bin/bookings ./cmd/web/*.go
./bin/bookings