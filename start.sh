#!/bin/bash

source .env

pmgo start "$(pwd)/bin/decision-maker" "$APP_NAME" true --args="-w=$(pwd)"
