package env

import (
	"github.com/eladhaim22/sic/utils"
)

var (
	ExcludeImages     = utils.GetEnv("EXCLUDE_IMAGES", "")
	CRON              = utils.GetEnv("CRON", "*/5 * * * *")
	SwarmMode         = utils.GetBooleanEnv("SWARM_MODE", false)
	NodesAgentsIp     = utils.GetEnv("NODES_AGENTS_IP", "")
	NodesAgentsPort   = utils.GetEnv("NODES_AGENTS_PORT", "")
)

