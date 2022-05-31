package log

import "github.com/rs/zerolog/log"

var Logger = log.With().Caller().Logger()
