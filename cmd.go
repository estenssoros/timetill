package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hako/durafmt"
	"github.com/spf13/cobra"
)

var (
	year     int
	month    int
	day      int
	hour     int
	minute   int
	second   int
	interval int
)

func init() {
	now := time.Now()
	cmd.Flags().IntVarP(&year, "year", "y", now.Year(), "year to countdown to")
	cmd.Flags().IntVarP(&month, "month", "M", int(now.Month()), "month to countdown to")
	cmd.Flags().IntVarP(&day, "day", "d", now.Day(), "day to countdown to")
	cmd.Flags().IntVarP(&hour, "hour", "H", now.Hour(), "hour to countdown to")
	cmd.Flags().IntVarP(&minute, "minute", "m", now.Minute(), "minute to countdown to")
	cmd.Flags().IntVarP(&second, "second", "s", now.Second(), "second to countdown to")
	cmd.Flags().IntVarP(&interval, "interval to send", "i", 1, "how frequently to send updates (minutes)")
}

func printDuration(deadline time.Time) {
	timeduration := deadline.Sub(time.Now())
	timeduration = timeduration.Round(time.Minute)
	duration := durafmt.Parse(timeduration)
	fmt.Println(duration)
}

var cmd = &cobra.Command{
	Use:   "timetill",
	Short: "countdown until a time",
	RunE: func(cmd *cobra.Command, args []string) error {
		deadline := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
		fmt.Printf("counting down until: %v\n", deadline)

		tickerCh := time.NewTicker(time.Duration(interval) * time.Minute).C
		stopChan := make(chan os.Signal, 2)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		for deadline.After(time.Now()) {
			printDuration(deadline)
			select {
			case <-tickerCh:
				printDuration(deadline)
			case <-stopChan:
				return nil
			}
		}
		fmt.Println("deadline passed!")
		return nil
	},
}

func execute() error {
	if err := cmd.Execute(); err != nil {
		return err
	}
	return nil
}
