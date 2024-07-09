package main

import "errors"

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")