package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/m1kx/go-vtr-backend/pkg/api"
	"github.com/m1kx/go-vtr-backend/pkg/config"
	"github.com/m1kx/go-vtr-backend/pkg/utils"
	"github.com/m1kx/go-vtr-backend/pkg/utils/health"
	"github.com/m1kx/go-vtr-backend/pkg/utils/notify"
	"github.com/m1kx/go-vtr-backend/pkg/utils/plan"
	"github.com/m1kx/go-vtr-backend/pkg/utils/pocketbase"
	"github.com/m1kx/go-vtr-backend/pkg/utils/structs"
)

func run(last_updated_at [2]string, last_num int) (new_updated_at [2]string, num_users int, err error) {
	scrape_start := time.Now()
	days := []string{"heute", "morgen"}
	users, err := pocketbase.GetAllUsers()
	if err != nil {
		return
	}
	scrape_taken := time.Since(scrape_start).Milliseconds()
	fmt.Printf("%sPocketBase%s user retrieve took %s%dms %s\n", config.Cyan, config.Reset, config.Cyan, scrape_taken, config.Reset)

	update_from_user := false
	verified_count := 0
	for _, obj := range users {
		if obj.VERIFIED {
			verified_count++
		}
		if obj.UPDATE {
			fmt.Printf("-> %sUser update%s from %s\n", config.Blue, config.Reset, obj.ID)
			update_from_user = true
			err = pocketbase.EditField("update", obj.ID, "users", false)
		}
	}
	num_users = verified_count

	for t := 0; t < len(days); t++ {
		scrape_site_start := time.Now()
		var data, updated_at, weekday, date_string, err = plan.Scrape(days[t])
		if err != nil || date_string == "" {
			fmt.Println(err)
			continue
		}
		go pocketbase.EditField(fmt.Sprintf("day_%d", t+1), "ux8ausqmf2h57dd", "times", date_string)

		scrape_site_taken := time.Since(scrape_site_start).Milliseconds()
		fmt.Printf("Site %sscrape%s for %s%s%s took %s%dms%s\n", config.Purple, config.Reset, config.Red, days[t], config.Reset, config.Purple, scrape_site_taken, config.Reset)

		check_start := time.Now()
		if updated_at == last_updated_at[t] && last_num == verified_count && !update_from_user {
			// website data didnt change since last time
			fmt.Println("Not running... [no data change]")
			new_updated_at = last_updated_at
			continue
		} else {
			last_updated_at[t] = updated_at
			new_updated_at[t] = updated_at
		}

		fmt.Println(fmt.Sprintf("Starting check for day %s...", days[t]))

		for i := 0; i < len(users); i++ {
			if !users[i].VERIFIED {
				continue
			}
			fmt.Println(fmt.Sprintf("Checking user %s", users[i].ID))
			subjects := strings.Split(users[i].SUBJECTS, ":")
			class := users[i].CLASS
			all := [][]string{}
			msg := fmt.Sprintf(">>> %s <<<", weekday)
			for x := 0; x < len(subjects); x++ {
				for d := 0; d < len(data); d++ {
					if data[d][0] == class && data[d][2] == subjects[x] {
						msg = fmt.Sprintf("%s\n%s", msg, notify.AssembleMessage(data, d))
						all = append(all, data[d])
					}
				}
			}

			day := ""
			if t == 0 {
				day = "h"
			} else {
				day = "m"
			}

			if len(all) == 0 {
				// clear hash
				hash_of_day := ""
				if day == "h" {
					hash_of_day = users[i].H_HASH
				} else {
					hash_of_day = users[i].M_HASH
				}
				if hash_of_day != "" {
					go pocketbase.EditField(fmt.Sprintf("%s_hash", day), users[i].ID, "users", "")
				}
				continue
			}
			all_string := ""
			all_eva := 0
			for i := 0; i < len(all); i++ {
				split := strings.Split(all[i][1], " - ")
				if strings.Contains(all[i][1], " - ") && len(split) == 2 && all[i][5] == "EVA" {
					from, _err := strconv.Atoi(split[0])
					to, _err := strconv.Atoi(split[1])
					if _err == nil {
						all_eva += -(from - to - 1)
					}
				} else if all[i][5] == "EVA" {
					all_eva += 1
				}
				for x := 0; x < len(all[i]); x++ {
					if x == len(all[i])-1 {
						all_string = fmt.Sprintf("%s%s!!!", all_string, all[i][x])
						continue
					}
					all_string = fmt.Sprintf("%s%s|", all_string, all[i][x])
				}
			}
			all_string = all_string[:len(all_string)-3]

			all_base := utils.EncodeBase64(all_string)
			if (day == "h" && all_base == users[i].H_HASH) || (day == "m" && all_base == users[i].M_HASH) {
				fmt.Printf("   -> %sNo update%s\n", config.Red, config.Reset)
				continue
			}

			if all_eva > 0 {
				go pocketbase.EditField("score", users[i].ID, "users", users[i].SCORE+all_eva)
			}

			if users[i].NEW_VERSION {
				notify.SendMail(msg, users[i].EMAIL)
			}

			if users[i].REQINFO.URL != "" && users[i].NEW_VERSION {
				props := ""
				url := users[i].REQINFO.URL
				if users[i].REQINFO.METHOD == "POST" {
					props = strings.Replace(users[i].REQINFO.INFOFMT, "TITLE", weekday, 1)
					props = strings.Replace(props, "MESSAGE", strings.Replace(msg, "\n", "\\n", -1), 1)
				} else if users[i].REQINFO.METHOD == "GET" {
					url = strings.Replace(url, "MESSAGE", strings.Replace(strings.Replace(msg, "\n", "%0A", -1), " ", "+", -1), 1)
				}
				data := structs.HttpReq{
					METHOD:   users[i].REQINFO.METHOD,
					PROPS:    props,
					BASE_URL: url,
				}
				err := notify.SendPerRequest(&data)
				if err != nil {
					fmt.Println(err)
				}
			}

			send_hash := ""
			send_hash = all_base
			go pocketbase.EditField(fmt.Sprintf("%s_hash", day), users[i].ID, "users", send_hash)
			//update_hash(send_hash, token, users[i].ID, day)

		}

		check_taken := time.Since(check_start).Milliseconds()
		fmt.Printf("%sData check%s on %s took %s%dms%s\n", config.Green, config.Reset, days[t], config.Green, check_taken, config.Reset)
	}
	return
}

func main() {

	godotenv.Load(".env")

	api.RunServer()

	args := os.Args[1:]

	last_updated_at := [2]string{"", ""}
	last_num := 0
	var err error

	// only for testing
	if len(args) > 0 && args[0] == "false" {
		for i := 0; i < 1000; i++ {
			last_updated_at, last_num, err = run(last_updated_at, last_num)
			if err != nil {
				fmt.Printf("Error occured:\n%s", err)
				err = nil
			}
			time.Sleep(3 * time.Second)
		}
		return
	}

	interval := time.Second
	currentTime := time.Now().Local()
	for {
		if (currentTime.Hour() == 7 && currentTime.Minute() >= 20) || (currentTime.Hour() == 8 && currentTime.Minute() <= 30) {
			interval = time.Minute
		} else if currentTime.Hour() >= 8 && currentTime.Hour() < 17 {
			interval = time.Minute * 5
		} else {
			interval = time.Minute * 20
		}
		last_updated_at, last_num, err = run(last_updated_at, last_num)
		if err != nil {
			fmt.Printf("Error occured:\n%s\n", err)
			health.Dead(err.Error())
			err = nil
			interval = time.Minute
		}
		time.Sleep(interval)
		currentTime = time.Now().Local()
	}

}
