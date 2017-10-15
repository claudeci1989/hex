package commands

import (
	"sort"
	"strings"

	"github.com/hexbotio/hex/models"
	"github.com/hexbotio/hex/parse"
)

func Help(message *models.Message, rules *map[string]models.Rule, config models.Config) {

	help := make([]string, 0)

	// pull internal help
	help = make([]string, 4)
	help[0] = "help <filter> - This help"
	help[1] = "ping - Simple ping response for the bot"
	help[2] = "rules - dump of loaded rules"
	help[3] = "version - Compiled version number/sha"

	// pull all help from rules
	for _, rule := range *rules {
		if rule.Active && !rule.Hide {
			if parse.Member(rule.ACL, message.Attributes["hex.user"]) || parse.Member(rule.ACL, message.Attributes["hex.channel"]) {
				if rule.Help != "" {
					help = append(help, rule.Help)
				} else {
					help = append(help, rule.Match)
				}
			}
		}
	}

	// sort, filter and de-dupe help
	helpTitle := "Help for " + config.BotName + "..."
	sort.Strings(help)
	var lasthelp = ""
	var newhelp = make([]string, 0)
	var filter = strings.TrimSpace(strings.Replace(message.Attributes["hex.input"], "help", "", 1))
	newhelp = append(newhelp, helpTitle)
	for _, h := range help {
		if strings.Contains(h, filter) || filter == "" {
			if h != lasthelp {
				newhelp = append(newhelp, " - "+h)
			}
			lasthelp = h
		}
	}

	// output help
	if len(newhelp) > 1 {
		message.Outputs = append(message.Outputs, models.Output{
			Rule:     "help",
			Response: strings.Join(newhelp[:], "\n"),
		})
	}

}
