package nodegraph_log

import "log"

var LOGGER log.Logger

func BuildLogger() {
	LOGGER = *log.New(log.Writer(), "[nodegraph plugin] ", log.LstdFlags|log.Lmsgprefix)
}
