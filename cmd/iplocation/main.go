package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/wzshiming/iplocation/server"
	"github.com/wzshiming/iplocation/source"
)

var (
	addr string
)

var (
	cmdRoot = cobra.Command{
		Use: "iplocation",
	}
	cmdServer = cobra.Command{
		Use:   "server",
		Short: "Run iplocation server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			server := server.NewServer(source.Sources)
			return http.ListenAndServe(addr, server.ServeMux())
		},
	}
	cmdQuery = cobra.Command{
		Use:   "query [source] [ip]",
		Short: "Run iplocation query",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, ok := source.Sources[args[0]]
			if !ok {
				return fmt.Errorf("source not found: %s", args[0])
			}
			addr, err := netip.ParseAddr(args[1])
			if err != nil {
				return err
			}
			data, err := src.Lookup(addr)
			if err != nil {
				return err
			}
			return json.NewEncoder(os.Stdout).Encode(data)
		},
	}
	cmdList = cobra.Command{
		Use:   "list",
		Short: "Run iplocation list source",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			keys := make([]string, 0, len(source.Sources))
			for k := range source.Sources {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Println(k)
			}
			return nil
		},
	}
)

func main() {
	ctx := context.Background()
	cmdServer.Flags().StringVar(&addr, "addr", ":8080", "http service address")
	cmdRoot.AddCommand(&cmdServer)
	cmdRoot.AddCommand(&cmdQuery)
	cmdRoot.AddCommand(&cmdList)
	err := cmdRoot.ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
