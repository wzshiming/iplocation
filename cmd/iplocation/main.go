package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wzshiming/iplocation/inject"
	"github.com/wzshiming/iplocation/server"
	"github.com/wzshiming/iplocation/source"
)

var (
	addr        string
	querySource string = "sapics-geolite2-city"
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
		Use:   "query [ip]",
		Short: "Run iplocation query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, ok := source.Sources[querySource]
			if !ok {
				return fmt.Errorf("source not found: %s", querySource)
			}
			addr, err := netip.ParseAddr(args[0])
			if err != nil {
				return err
			}
			data, err := src.Lookup(addr)
			if err != nil {
				data = map[string]any{
					"error": err.Error(),
				}
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(data)
		},
	}
	injectServer = cobra.Command{
		Use:   "inject",
		Short: "Inject IP location",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			src, ok := source.Sources[querySource]
			if !ok {
				return fmt.Errorf("source not found: %s", querySource)
			}
			_, err := io.Copy(os.Stdout, inject.NewReader(os.Stdin, func(addr netip.Addr) ([]byte, bool) {
				data, err := src.Lookup(addr)
				if err != nil {
					return []byte(fmt.Sprintf("<>%s", addr)), true
				}
				return []byte(fmt.Sprintf("<%s>%s", data, addr)), true
			}))
			return err
		},
	}
	cmdList = cobra.Command{
		Use:   "list",
		Short: "Run iplocation list source",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(allSourceList)
		},
	}
)

var allSourceList = allSource()

func allSource() []string {
	keys := make([]string, 0, len(source.Sources))
	for k := range source.Sources {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func main() {
	ctx := context.Background()
	cmdServer.Flags().StringVar(&addr, "addr", ":8080", "http service address")
	cmdRoot.AddCommand(&cmdServer)
	cmdRoot.AddCommand(&cmdQuery)
	cmdRoot.AddCommand(&cmdList)
	cmdRoot.AddCommand(&injectServer)
	cmdRoot.PersistentFlags().StringVar(&querySource, "source", querySource, "one of ("+strings.Join(allSourceList, ", ")+")")
	err := cmdRoot.ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
