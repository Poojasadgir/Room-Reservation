#!/bin/bash

go build -o room-reservation cmd/web/*.go && ./room-reservation -dbname=postgres -dbuser=postgres -cache=false -production=false