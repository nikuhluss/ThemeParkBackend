package cmd

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/generator"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates and populates the database with mock data",
	Run: func(cmd *cobra.Command, args []string) {

		var db *sqlx.DB
		var err error

		dbconfig := testutil.NewDatabaseConnectionConfig()

		if dokku {
			db, err = testutil.NewDatabaseConnectionDokku(dbconfig)
		} else {
			db, err = testutil.NewDatabaseConnection(dbconfig)
		}

		if err != nil {
			fmt.Printf("error creating db connection: %s\n", err)
			os.Exit(1)
		}

		// begin transaction

		tx, err := db.Begin()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		inserter := generator.NewInserter(tx)
		inserter.Seed(1)
		err = inserter.DoInsert()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// end transaction

		err = tx.Commit()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
