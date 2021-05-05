package env

import (
	"github.com/eladhaim22/sic/utils"
)

var (
	ExcludeImages              = utils.Getenv("EXCLUDE_IMAGES", "")
	CRON                       = utils.Getenv("CRON", "*/30 * * * * *")
)

