package cmd

import (
	"log"
	"time"

	"github.com/hikaru7719/concat-gzip-log/pkg/runner"
	"github.com/spf13/cobra"
)

type Command struct {
	*cobra.Command
}

func NewCommand() *Command {
	command := &cobra.Command{
		Use:   "concat-gzip-log",
		Short: "A concat cli for gzip access log file on aws s3",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := cmd.Flags().GetString("bucket")
			if err != nil {
				log.Print(err)
				return err
			}

			d, err := cmd.Flags().GetString("date")
			if err != nil {
				log.Print(err)
				return err
			}

			n, err := cmd.Flags().GetString("name")
			if err != nil {
				log.Print(err)
				return err
			}

			t, err := convertToTime(d)
			if err != nil {
				log.Print(err)
				return err
			}

			runner.NewRunner(b, t, n).Run()
			return nil
		},
	}
	command.PersistentFlags().StringP("bucket", "b", "", "specify target bucket name")
	command.PersistentFlags().StringP("date", "d", "", "specify target date YYYY-MM-DD")
	command.PersistentFlags().StringP("name", "n", "concat-gzip-log.txt", "specify output file name")
	command.PersistentFlags().BoolP("parallel", "p", false, "run parallel")
	return &Command{command}
}

func convertToTime(date string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
