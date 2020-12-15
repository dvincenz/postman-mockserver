
package main

import (
	"github.com/dvincenz/postman-mockserver/cmd"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	cmd.Execute()

}
