#!/bin/bash
curl -v -X POST 'localhost:3030/session' -d '{"data":{"type":"session","id":"poop","attributes":{"created":"2015-05-12T20:40:02.291562383-05:00"},"relationships":{"logged_in_as":{"linkage":[{"id":"3","type":"user"}]}}}}'
curl 'localhost:3030/session/poop'
curl -v -X POST 'localhost:3030/session' -d '{"data":{"type":"session","id":"poop","attributes":{"created":"2015-05-12T20:40:02.291562383-05:00"}}}'
