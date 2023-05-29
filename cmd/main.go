package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rafata1/feedies/config"
	"github.com/rafata1/feedies/handler/feed_handler"
	"github.com/rafata1/feedies/service/feed_service"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"time"
)

const (
	versionTimeFormat = "20060102150405"
)

func main() {
	cmd := &cobra.Command{}
	cmd.AddCommand(
		createMigrationCommand(),
		migrateCommand(),
		consumeRSSCommand(),
		startServerCommand(),
	)
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

func createMigrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-create [name]",
		Short: "create sql migrations",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			version := now.Format(versionTimeFormat)
			name := args[0]
			up := fmt.Sprintf("%s/%s_%s.up.sql", "migration", version, name)
			down := fmt.Sprintf("%s/%s_%s.down.sql", "migration", version, name)

			err := os.WriteFile(up, []byte{}, 0644)
			if err != nil {
				panic(err)
			}

			err = os.WriteFile(down, []byte{}, 0644)
			if err != nil {
				panic(err)
			}

			fmt.Println("Created SQL up script:", up)
			fmt.Println("Created SQL down script:", down)
		},
	}
}

func migrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-up",
		Short: "migrate all the way up",
		Run: func(cmd *cobra.Command, args []string) {
			conf := config.Load()
			m, err := migrate.New(
				fmt.Sprintf("file://%s", "migration"),
				fmt.Sprintf("mysql://%s", conf.MySQL.DSN()),
			)
			if err != nil {
				panic(err)
			}

			err = m.Up()
			if err == migrate.ErrNoChange {
				fmt.Println("No change in migration")
				return
			}
			if err != nil {
				panic(err)
			}
			fmt.Println("Migrated up")
		},
	}
}

func consumeRSSCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "consume-rss",
		Short: "consume rss from sources",
		Run: func(cmd *cobra.Command, args []string) {
			conf := config.Load()
			db := conf.MySQL.MustConnect()
			feedService := feed_service.Init(db)
			feedService.ConsumeRSS(context.Background())
		},
	}
}

func startServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "start server",
		Run: func(cmd *cobra.Command, args []string) {
			conf := config.Load()
			feedHandler := feed_handler.Init(conf)

			router := gin.Default()
			router.Use(cors.Default())
			router.GET("/", func(c *gin.Context) {
				c.String(http.StatusOK, "Hello World!")
			})
			router.GET("/api/v1/news", feedHandler.GetNews)
			router.Run()
		},
	}
}
