package duration

import (
	"fmt"
	"github.com/k8spacket/plugins/nodegraph/metrics/nodegraph/model"
	"time"
)

func GetConfig() model.Config {
	return model.Config{Arc1: model.DisplayConfig{DisplayName: "Average duration", Color: "purple"},
		Arc2:          model.DisplayConfig{DisplayName: "Max duration", Color: "white"},
		MainStat:      model.DisplayConfig{DisplayName: "Average duration "},
		SecondaryStat: model.DisplayConfig{DisplayName: "Max duration "}}
}

func FillNodeStats(node *model.Node, connEndpoint model.ConnectionEndpoint) {
	if connEndpoint.Duration > 0 || connEndpoint.MaxDuration > 0 {
		var cd = connEndpoint.Duration / float64(connEndpoint.ConnCount)
		if cd >= 0.001 {
			node.MainStat = fmt.Sprintf("avg: %s", time.Duration(cd*float64(time.Millisecond)))
		} else {
			node.MainStat = fmt.Sprint("avg: <0.001s")
		}
		if connEndpoint.MaxDuration >= 0.001 {
			node.SecondaryStat = fmt.Sprintf("max: %s", time.Duration(connEndpoint.MaxDuration*float64(time.Millisecond)))
		} else {
			node.SecondaryStat = fmt.Sprint("max: <0.001s")
		}
		node.Arc1 = cd / connEndpoint.MaxDuration
		node.Arc2 = (connEndpoint.MaxDuration - cd) / connEndpoint.MaxDuration
	} else {
		node.MainStat = fmt.Sprint("avg: N/A")
		node.SecondaryStat = fmt.Sprint("max: N/A")
	}
}

func FillEdgeStats(edge *model.Edge, connItem model.ConnectionItem) {
	if connItem.Duration > 0 || connItem.MaxDuration > 0 {
		var cd = connItem.Duration / float64(connItem.ConnCount)
		if cd >= 0.001 {
			edge.MainStat = fmt.Sprintf("avg: %s", time.Duration(cd*float64(time.Millisecond)))
		} else {
			edge.MainStat = fmt.Sprint("avg: <0.001s")
		}
		if connItem.MaxDuration >= 0.001 {
			edge.SecondaryStat = fmt.Sprintf("max: %s", time.Duration(connItem.MaxDuration*float64(time.Millisecond)))
		} else {
			edge.SecondaryStat = fmt.Sprint("max: <0.001s")
		}
	} else {
		edge.MainStat = fmt.Sprint("avg: N/A")
		edge.SecondaryStat = fmt.Sprint("max: N/A")
	}
}
