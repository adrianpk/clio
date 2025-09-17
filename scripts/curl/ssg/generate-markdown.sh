#!/bin/bash
# Triggers the site generation process.

curl -i -X POST http://localhost:8081/api/v1/ssg/generate-markdown
