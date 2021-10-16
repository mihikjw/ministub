#!/bin/bash
APP_NAME=ministub

go clean
if [ -f bin/${APP_NAME} ]; then rm bin/${APP_NAME}; fi;
if [ -f .coverage.html ]; then rm .coverage.html; fi;
if [ -f coverage.out ]; then rm coverage.out; fi;