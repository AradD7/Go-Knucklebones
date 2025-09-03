#!/bin/bash

curl -X POST\
    -H "Content-Type: application/json"\
    -d '{"username": "golerp", "password": "golepJuice"}'\
    -w "\nStatus: %{http_code}\n"\
    "http://localhost:8080/api/players"

