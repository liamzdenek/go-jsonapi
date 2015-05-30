#!/bin/bash
curl 'localhost:3030/session/1?include=logged_in_as,logged_in_as.posts,logged_in_as.posts.comments' | python -m json.tool
