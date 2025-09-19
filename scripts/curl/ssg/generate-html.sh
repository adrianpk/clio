#!/bin/bash
# Triggers the site HTML generation process.

curl -i -X POST http://localhost:8081/api/v1/ssg/generate-html
