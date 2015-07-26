#!/bin/bash
curl -v -X POST 'localhost:3030/user' -d '{"data":{"type":"user","attributes":{"name":"testy mctesterson"},"relationships":{"logged_in_as":{"data":[{"id":"3","type":"user"}]}}}}'
