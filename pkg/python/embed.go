package python

import _ "embed"

//go:embed create.py
var CreateScript string

//go:embed startup.py
var StartupScript string
