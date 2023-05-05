package plan

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func format_row(row []string) {
	split := strings.Split(row[2], " ")
	if len(split) > 2 {
		row[2] = fmt.Sprintf("%s %s", split[0], split[2])
	}
}

func Scrape(day string) (data [][]string, updated_at string, wd string, date_string string, err error) {
	res, err := http.Get(fmt.Sprintf("https://lmg-anrath.de/aktuelle_plaene/Vertretungsplan/%s/subst_001.htm", day))
	if err != nil {
		return nil, "", "", "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, "", "", "", err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", "", "", err
	}

	layout := "Mon, 02 Jan 2006 15:04:05 MST"
	updated_at = res.Header.Get("last-modified")
	updated_date, err := time.Parse(layout, updated_at)
	if err != nil {
		fmt.Println(err)
	}

	// date handling
	heading := doc.Find(".mon_title").Text()
	heading_fmt := strings.ReplaceAll(strings.Split(heading, " ")[0], ".", "-")
	date, err := time.Parse("2-1-2006", heading_fmt)
	if err != nil {
		return nil, "", "", "", err
	}
	date_string = date.Format("02-01-2006")
	curr_time := time.Now()
	if day == "morgen" {
		curr_time = curr_time.AddDate(0, 0, 1)
	}
	curr_date := curr_time.Format("02-01-2006")
	_ = curr_date

	weekday := ""
	switch date.Weekday() {
	case time.Sunday:
		weekday = "Sonntag"
	case time.Monday:
		weekday = "Montag"
	case time.Tuesday:
		weekday = "Dienstag"
	case time.Wednesday:
		weekday = "Mittwoch"
	case time.Thursday:
		weekday = "Donnerstag"
	case time.Friday:
		weekday = "Freitag"
	case time.Saturday:
		weekday = "Samstag"
	}

	if curr_time.Weekday() == date.Weekday() {
		weekday = fmt.Sprintf("%s (%s)", weekday, day)
	}

	var rows []*goquery.Selection
	doc.Find(".list").Each(func(i int, s *goquery.Selection) {
		rows = append(rows, s)
	})

	if strings.Contains(heading, "Seite") {
		sites := strings.Split(strings.Split(heading, "(Seite")[1], " / ")
		min_site, err := strconv.Atoi(strings.TrimLeft(sites[0], " "))
		max_site, err := strconv.Atoi(strings.Trim(sites[1], ")"))
		if err == nil {
			for i := min_site + 1; i <= max_site; i++ {
				fmt.Printf("Scraping multiple sites... %d\n", i)
				res, err := http.Get(fmt.Sprintf("https://lmg-anrath.de/aktuelle_plaene/Vertretungsplan/%s/subst_00%d.htm", day, i))
				if err == nil {
					updated_at_new := res.Header.Get("last-modified")
					updated_date_new, err := time.Parse(layout, updated_at_new)
					if err != nil {
						fmt.Println(err)
					} else {
						if updated_date.Before(updated_date_new) {
							updated_at = updated_at_new
						}
					}
					doc, err := goquery.NewDocumentFromReader(res.Body)
					if err == nil {
						doc.Find(".list").Each(func(x int, s *goquery.Selection) {
							if x > 7 {
								rows = append(rows, s)
							}
						})
					}
				}
			}
		}
	} else {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrHandlerTimeout
			},
		}
		for i := 2; i < 5; i++ {
			res, err := client.Get(fmt.Sprintf("https://lmg-anrath.de/aktuelle_plaene/Vertretungsplan/%s/subst_00%d.htm", day, i))
			if err != nil {
				continue
			}
			
			updated_at_new := res.Header.Get("last-modified")
			updated_date_new, err := time.Parse(layout, updated_at_new)
			if err != nil {
				fmt.Println(err)
			} else {
				if updated_date.Before(updated_date_new) {
					updated_at = updated_at_new
				}
			}
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err == nil {
				heading_new := doc.Find(".mon_title").Text()
				heading_fmt_new := strings.ReplaceAll(strings.Split(heading_new, " ")[0], ".", "-")
				date_new, err := time.Parse("2-1-2006", heading_fmt_new)
				if err != nil {
					fmt.Println(err)
				}
				date_string_new = date_new.Format("02-01-2006")
				if date_string_new != date_string {
					break
				}
				doc.Find(".list").Each(func(x int, s *goquery.Selection) {
					if x > 7 {
						rows = append(rows, s)
					}
				})
			}
		}
	}

	var rows_formatted [][]string
	for i := 9; i < len(rows); i++ {
		var this_row []string
		for x := 0; x < 7; x++ {
			this_row = append(this_row, rows[i+x].Text())
		}
		i += 7
		rows_formatted = append(rows_formatted, this_row)
	}

	for i := 0; i < len(rows_formatted); i++ {
		format_row(rows_formatted[i])
	}

	return rows_formatted, updated_at, weekday, date_string, nil
}
