package ws

import "go.aew.app/api.v1/log"

var (
	cPool map[string]*WSCtx = make(map[string]*WSCtx)
)

func setWs(label string, ws *WSCtx) {
	if v, b := cPool[label]; b {
		log.Infof("Exist Ctx %s", v.label)
	} else {
		cPool[label] = ws
	}
}

func GetCtx(label string) *WSCtx {
	if v1, bExit := cPool[label]; bExit {
		return v1
	}
	return nil
}
