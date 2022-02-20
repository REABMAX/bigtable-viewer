# Bigtable Viewer - This is currently a WIP

Currently a WIP without any tests ;)

## Getting started

    docker-compose up -d
    make init
    make run

Then, visit `localhost:3001`

## Options

Options can be passed by either defining an ENV var or adding a flag to the cli execution.

**-project (or env: PROJECT)**  
The gcp project id

**-instance (or env: INSTANCE)**  
Your bigtable instance

## Using locally with the emulator

Use with the env var `BIGTABLE_EMULATOR_HOST=localhost:8086`