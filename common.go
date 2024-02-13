package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func loadItemsFromFile(filename string) ([]*Item, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	items := make([]*Item, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " ")
		if len(fields) >= 4 {
			name := fields[0]
			link := fields[1]
			tags := fields[2]
			dateStr := fields[3]

			date, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				return nil, err
			}

			item := &Item{
				Name: name,
				Date: date,
				Tags: tags,
				Link: link,
			}

			items = append(items, item)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func saveItemsToFile(items []*Item, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range items {
		line := fmt.Sprintf("%s %s %s %s\n", item.Name, item.Link, item.Tags, item.Date.Format(time.RFC3339))
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}
