package core

var KeySpaceStat [4]map[string]int

func init() {
	// each db slot needs its own map, otherwise writes panic on a nil map
	for i := range KeySpaceStat {
		KeySpaceStat[i] = make(map[string]int)
	}
}

func UpdateKeySpaceStat(num int, metric string, val int) {
	KeySpaceStat[num][metric] = val
}
