package notify

import (
	"fmt"
	"strings"
)

func AssembleMessage(data [][]string, d int) (msg_part string) {
	info_type := data[d][5]
	info_room := ""
	if strings.Contains(info_type, "Raum") {
		info_room = fmt.Sprintf(" in Raum %s", data[d][4])
	} else if strings.Contains(info_type, "Klausur") {
		info_room = fmt.Sprintf(" in Raum %s", data[d][4])
	}
	info_time := data[d][1]
	info_time_seperated := strings.Split(info_time, " - ")
	if len(info_time_seperated) > 1 {
		info_time = fmt.Sprintf(" von Stunde %s bis %s ", info_time_seperated[0], info_time_seperated[1])
	} else {
		info_time = fmt.Sprintf(" in Stunde %s ", info_time)
	}
	info_subject := strings.Split(data[d][3], " ")[0]
	msg_part = fmt.Sprintf("🤖 %s in %s%s%s", info_type, info_subject, info_room, info_time)
	return
}
