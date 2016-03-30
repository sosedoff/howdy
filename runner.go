package main

import (
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"
)

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func RunCheck(config *Config, check Check) string {
	ts := time.Now()
	err := check.Perform()
	duration := Round(time.Since(ts).Seconds()*1000, .5, 2)

	record := CheckRecord{
		RunId:     RunId,
		Service:   config.Name,
		Check:     strings.ToLower(check.Name()),
		Success:   err == nil,
		Duration:  duration,
		Timestamp: ts,
	}

	if err != nil {
		record.Message = err.Error()
	}

	DB.Create(&record)

	if err != nil {
		log.Println("Check failed:", err)
		return fmt.Sprintf(
			"%v check failed for %v: %v",
			check.Name(), config.Name, err.Error(),
		)
	}

	return ""
}

func RunConfig(config *Config, wg *sync.WaitGroup) {
	messages := []string{}

	for _, check := range config.Checks {
		msg := RunCheck(config, check)
		if msg != "" {
			messages = append(messages, msg)
		}
	}

	if SendNotifications {
		for _, notifier := range config.Notifiers {
			err := notifier.Perform(messages)
			if err != nil {
				log.Println("Notifier failed:", err)
			}
		}
	}

	if wg != nil {
		wg.Done()
	}
}
