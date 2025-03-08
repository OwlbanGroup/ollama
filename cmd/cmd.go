package cmd

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"github.com/sirupsen/logrus"

	"math"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/containerd/console"
	"github.com/mattn/go-runewidth"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"github.com/ollama/ollama/api"
	"github.com/ollama/ollama/auth"
	"github.com/ollama/ollama/envconfig"
	"github.com/ollama/ollama/format"
	"github.com/ollama/ollama/parser"
	"github.com/ollama/ollama/progress"
	"github.com/ollama/ollama/server"

	"github.com/ollama/ollama/types/errtypes"
	"github.com/ollama/ollama/types/model"
	"github.com/ollama/ollama/version"
)

/*
 * CreateHandler handles the creation of a model from a Modelfile.
 * It reads the specified file, processes the commands, and interacts with the API client.
 */
	if err != nil {
		logrus.Error(err)
		return err

	client, err := api.ClientFromEnvironment()
	if err != nil {
		logrus.Error(err)
		return err
	}

	p := progress.NewProgress(os.Stderr)
	defer p.Stop()

	f, err := os.Open(filename)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer f.Close()

	modelfile, err := parser.ParseFile(f)
	if err != nil {
		logrus.Error(err)
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Error(err)
		return err
	}

	status := "transferring model data"
	spinner := progress.NewSpinner(status)
	p.Add(status, spinner)

	bars := make(map[string]*progress.Bar)

	for i := range modelfile.Commands {
		switch modelfile.Commands[i].Name {
		case "model", "adapter":
			path := modelfile.Commands[i].Args
			if path == "~" {
				path = home
			} else if strings.HasPrefix(path, "~/") {
				path = filepath.Join(home, path[2:])
			}

			if !filepath.IsAbs(path) {
				path = filepath.Join(filepath.Dir(filename), path)
			}

			fi, err := os.Stat(path)
			if errors.Is(err, os.ErrNotExist) && modelfile.Commands[i].Name == "model" {
				continue
			} else if err != nil {
				logrus.Error(err)
				return err
			}

			if fi.IsDir() {
				tempfile, err := tempZipFiles(path)
				if err != nil {
					logrus.Error(err)
					return err
				}
				defer os.RemoveAll(tempfile)

				path = tempfile
			}

			digest, err := createBlob(cmd, client, path)
			if err != nil {
				logrus.Error(err)
				return err
			}

			modelfile.Commands[i].Args = "@" + digest
		}
	}

	fn := func(resp api.ProgressResponse) error {
		if resp.Digest != "" {
			spinner.Stop()

			bar, ok := bars[resp.Digest]
			if !ok {
				bar = progress.NewBar(fmt.Sprintf("pulling %s...", resp.Digest[7:19]), resp.Total, resp.Completed)
				bars[resp.Digest] = bar
				p.Add(resp.Digest, bar)
			}

			bar.Set(resp.Completed)
		} else if status != resp.Status {
			spinner.Stop()

			status = resp.Status
			spinner = progress.NewSpinner(status)
			p.Add(status, spinner)
		}

		return nil
	}

	request := api.CreateRequest{Name: args[0], Modelfile: modelfile.String()}

	if err := client.Create(cmd.Context(), &request, fn); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

// Other functions remain unchanged...
