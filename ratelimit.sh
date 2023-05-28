#!/bin/bash

for i in {1..6}; do curl http://localhost:4000/v1/healthcheck; done

