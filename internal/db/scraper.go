package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"
)

// get code recency options
type CodeRecencyOption uint8
const AllCodes CodeRecencyOption = 0
const RecentCodes CodeRecencyOption = 1
const UnrecentCodes CodeRecencyOption = 2

func AddCode(code string, game string, description string, livestream bool, foundTime time.Time) error {
	_, err := DBScraper.Exec("INSERT INTO Codes SET code = ?, game = ?, description = ?, is_livestream = ?, added = ?", code, game, description, livestream, foundTime)
	return err
}

func RemoveCode(code string) error {
	_, err := DBScraper.Exec("DELETE FROM FROM Codes WHERE code = ?", code)
	return err
}

func GetMostRecentCodeTime(game string) (time.Time, error) {
	var time time.Time
	sel := DBScraper.QueryRow("SELECT added FROM Codes WHERE game = ? ORDER BY added DESC", game)
	err := sel.Scan(&time)
	return time, err
}

func GetCodes(game string, recency CodeRecencyOption, livestream bool) [][]string {
	var sels *sql.Rows
	var err error
	codes := [][]string{}

	switch recency {
	case AllCodes:
		sels, err = DBScraper.Query("SELECT code, description FROM Codes WHERE game = ? AND is_livestream = ? ORDER BY added ASC", game, livestream)
	case RecentCodes:
		// get most recent code's added datetime
		recentTime, rerr := GetMostRecentCodeTime(game)
		if rerr != nil {
			log.Fatalf("Error getting most recent code time for %v: %v", game, rerr)
		}
		// get codes added within 24 hours before the most recent
		oldestTime := recentTime.Add(-24 * time.Hour)
		sels, err = DBScraper.Query("SELECT code, description FROM Codes WHERE game = ? AND is_livestream = ? AND added >= ? ORDER BY added ASC", game, livestream, oldestTime)
	case UnrecentCodes:
		// get most recent code's added datetime
		recentTime, rerr := GetMostRecentCodeTime(game)
		if rerr != nil {
			log.Fatalf("Error getting most recent code time for %v: %v", game, err)
		}
		// select codes added older than 24 hours before the most recent
		oldestTime := recentTime.Add(-24 * time.Hour)
		sels, err = DBScraper.Query("SELECT code, description FROM Codes WHERE game = ? AND is_livestream = ? AND added < ? ORDER BY added ASC", game, livestream, oldestTime)
	}
	
	if err != nil {
		log.Fatalf("Error querying codes of recency %v: %v", recency, err)
	}

	var code string
	var description string
	for sels.Next() {
		sels.Scan(&code, &description)
		codes = append(codes, []string{code, description})
	}
	if err = sels.Err(); err != nil {
		log.Fatalf("Error reading code row for %v: %v", game, err)
	}

	return codes
}

func GetRemovedCodes(codes []string, game string, removeFromDB bool) ([][]string, error) {
	result := [][]string{}
	codesPlaceholder := Placeholders(len(codes))

	// convert to any slice
	queryArgs := make([]any, len(codes) + 1)
	queryArgs[0] = game
	for i, v := range codes {
		queryArgs[i+1] = v
	}

	q := fmt.Sprintf("SELECT code, description FROM Codes WHERE game = ? AND code NOT IN (%s)", codesPlaceholder)
	sels, err := DBScraper.Query(q, queryArgs...)
	if err != nil {
		return result, err
	}
	
	for sels.Next() {
		var code, desc string
		sels.Scan(&code, &desc)
		fmt.Printf("Found removed code: %v\n", code)
		result = append(result, []string{code, desc})
	}
	if sels.Err() != nil {
		return result, sels.Err()
	}

	// needed as otherwise would cause an error on delete command below
	if len(result) == 0 {
		return result, err
	}

	if removeFromDB {
		deleteArgs := make([]any, len(result) + 1)
		deleteArgs[0] = game
		for i, v := range result {
			deleteArgs[i+1] = v[0]
		}

		q := fmt.Sprintf("DELETE FROM Codes WHERE game = ? AND code IN (%s)", Placeholders(len(result)))
		_, err = DBScraper.Exec(q, deleteArgs...)
	}

	return result, err
}

func SetScrapeTimes(game string, updated time.Time, checked time.Time) error {
	row := DBScraper.QueryRow("SELECT game FROM ScrapeStats WHERE game = ?", game)
	
	var z string // temp unused var for existence checking
	err := row.Scan(&z)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info(fmt.Sprintf("Adding %v to ScrapeStats", game))
			_, err := DBScraper.Exec("INSERT INTO ScrapeStats SET game = ?", game)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	_, err = DBScraper.Exec("UPDATE ScrapeStats SET updated = ?, checked = ? WHERE game = ?", updated, checked, game)
	return err
}

// Returns time scraped, time source updated, and db read error.
func GetScrapeTimes(game string) (time.Time, time.Time, error) {
	var checked time.Time
	var updated time.Time
	row := DBScraper.QueryRow("SELECT checked, updated FROM ScrapeStats WHERE game = ?", game)
	err := row.Scan(&checked, &updated)
	return checked, updated, err
}