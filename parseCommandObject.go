package main

// CommandObject as name
type CommandObject struct {
	app            string
	flashver       string
	swfURL         string
	tcURL          string
	fpad           bool
	audioCodecs    float64
	videoCodecs    float64
	videoFunction  float64
	pageURL        string
	objectEncoding float64
}

func parseCommandObject(raw map[string]interface{}) CommandObject {
	cmdOjbect := CommandObject{}

	tmp, ok := raw["app"]
	if ok {
		tmp, ok := tmp.(string)
		if ok {
			cmdOjbect.app = tmp
		}
	}

	tmp, ok = raw["flashver"]
	if ok {
		tmp, ok := tmp.(string)
		if ok {
			cmdOjbect.flashver = tmp
		}
	}

	tmp, ok = raw["swfUrl"]
	if ok {
		tmp, ok := tmp.(string)
		if ok {
			cmdOjbect.swfURL = tmp
		}
	}

	tmp, ok = raw["tcUrl"]
	if ok {
		tmp, ok := tmp.(string)
		if ok {
			cmdOjbect.tcURL = tmp
		}
	}

	tmp, ok = raw["fpad"]
	if ok {
		tmp, ok := tmp.(bool)
		if ok {
			cmdOjbect.fpad = tmp
		}
	}

	tmp, ok = raw["audioCodecs"]
	if ok {
		tmp, ok := tmp.(float64)
		if ok {
			cmdOjbect.audioCodecs = tmp
		}
	}

	tmp, ok = raw["videoCodecs"]
	if ok {
		tmp, ok := tmp.(float64)
		if ok {
			cmdOjbect.videoCodecs = tmp
		}
	}

	tmp, ok = raw["videoFunction"]
	if ok {
		tmp, ok := tmp.(float64)
		if ok {
			cmdOjbect.videoFunction = tmp
		}
	}

	tmp, ok = raw["pageUrl"]
	if ok {
		tmp, ok := tmp.(string)
		if ok {
			cmdOjbect.pageURL = tmp
		}
	}

	tmp, ok = raw["objectEncoding"]
	if ok {
		tmp, ok := tmp.(float64)
		if ok {
			cmdOjbect.objectEncoding = tmp
		}
	}

	return cmdOjbect
}
