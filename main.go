package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/chartmuseum/helm-push/pkg/helm"
	"github.com/spf13/cobra"
)

var version = ""

func main() {
	cmd := &cobra.Command{
		Use:          "helm push [chart archive/directory] [repository name]",
		Short:        "push chart package to Coding artifact",
		Long:         "push chart package to Coding artifact",
		SilenceUsage: false,
		Args:         cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("chart archive/directory and repository name required")
			}
			return push(args[0], args[1])
		},
	}

	flags := cmd.Flags()
	flags.Parse(os.Args[1:])

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func push(chartName, repository string) error {
	chart, err := helm.GetChartByName(chartName)
	if err != nil {
		return err
	}

	repo, err := helm.GetRepoByName(repository)
	if err != nil {
		return err
	}

	tmp, err := ioutil.TempDir("", "helm-push-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	chartPackagePath, err := helm.CreateChartPackage(chart, tmp)
	if err != nil {
		return err
	}

	chartArchiveFile, err := os.Open(chartPackagePath)
	if err != nil {
		return err
	}
	defer chartArchiveFile.Close()

	fileInfo, err := chartArchiveFile.Stat()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", repo.URL, chartArchiveFile)
	if err != nil {
		return err
	}

	req.SetBasicAuth(repo.Username, repo.Password)
	req.Header.Set("User-Agent", fmt.Sprintf("helm-push/%s", version))
	req.Header.Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	fmt.Printf("pushing chart '%s' to repository '%s' ...\n", chart.Metadata.Name, repo.URL)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf(
			"fail to push chart '%s' to repository '%s': repository returns error %s: %s",
			chart.Metadata.Name,
			repo.URL,
			http.StatusText(resp.StatusCode),
			string(body),
		)
	}

	fmt.Printf("finished pushing chart '%s' to repository '%s'\n", chart.Metadata.Name, repo.URL)
	return nil
}
