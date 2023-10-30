package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
)

type crawlerServer struct {
	output    string
	pgConnUrl string // pg connection url
}

var (
	crawlerCmd = &cobra.Command{
		Use:   "crawl",
		Short: "crawl",
		Run: func(cmd *cobra.Command, args []string) {
			if err := crawler.start(); err != nil {
				log.Fatalln(err)
			}
			log.Println("succeeded")
		},
	}

	crawler = &crawlerServer{}
)

func (s *crawlerServer) validate() error {
	if s.pgConnUrl == "" {
		return fmt.Errorf("please specify postgres connection url with --postgres")
	}

	return nil
}

func (s *crawlerServer) start() error {
	if err := s.validate(); err != nil {
		log.Fatalln(err)
	}

	log.Println("start to crawl into file ", s.output)
	if err := s.crawlTenants(); err != nil {
		return err
	}

	return nil
}

func (s *crawlerServer) crawlTenants() error {
	conn, err := s.getDbConn()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	data := &tenant{}
	if err := conn.QueryRow(context.Background(),
		//`select "id","createdOn","displayName" from "Tenant" where "suite" = 'TENCENT_MEETING' and "status" = 'ACTIVE'`
		`select "id","createdOn","displayName" from "Tenant" where "id" = $1`, "tn-598738aac4d04a688e3583ace317e57c",
	).Scan(data.ID, data.CreatedOn, data.DisplayName); err != nil {
		return err
	}

	log.Println(data.String())
	return nil
}

// get db connection
func (s *crawlerServer) getDbConn() (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), s.pgConnUrl)
}

func (s *crawlerServer) defaultOutput() string {
	return "./reports/" + time.Now().Format("2006-01-02")
}

func init() {
	crawlerCmd.PersistentFlags().StringVarP(&crawler.output,
		"output", "o", crawler.defaultOutput(), "output report dir",
	)
	crawlerCmd.PersistentFlags().StringVarP(&crawler.pgConnUrl,
		"postgres", "u", "", "postgres connection url, for example postgres://username:password@localhost:5432/database_name",
	)
	crawlerCmd.MarkFlagRequired("postgres")
	rootCmd.AddCommand(crawlerCmd)
}

type (
	tenant struct {
		ID          string
		DisplayName string
		CreatedOn   int64
	}
)

func (t tenant) String() string {
	return fmt.Sprintf("id=%s displayName=%s createdOn=%v", t.ID, t.DisplayName, t.CreatedOn)
}
