package client

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "embed"

	"git.sr.ht/~mpldr/uniview/internal/buildinfo"
	"git.sr.ht/~mpldr/uniview/internal/client/api"
	"git.sr.ht/~mpldr/uniview/internal/config"
	"git.sr.ht/~mpldr/uniview/internal/player"
	"github.com/ogen-go/ogen/ogenerrors"
)

func StartRestServer(_ context.Context, p player.Interface, status *api.StatusConnection, remote *url.URL) error {
	r := &restServer{p, status}

	srv, err := api.NewServer(r, api.WithErrorHandler(ogenerrors.DefaultErrorHandler))
	if err != nil {
		return fmt.Errorf("failed to create client API: %w", err)
	}

	wrapper := NewDocWrapper(srv, remote)

	httpSrv := &http.Server{
		Addr:    "[::1]:21558",
		Handler: wrapper,
	}

	shutdownFuncs = append(shutdownFuncs, httpSrv.Shutdown)

	return httpSrv.ListenAndServe()
}

type restServer struct {
	p      player.Interface
	status *api.StatusConnection
}

// PlayerStartPost implements POST /player/start operation.
// Start playback of a video.
//
// POST /player/start
func (r *restServer) PlayerStartPost(ctx context.Context, req api.OptPlayerStartPostReq) (api.PlayerStartPostRes, error) {
	openPath := ""
	switch req.Value.Type {
	case api.VideoFilePlayerStartPostReq:
		file, ok := req.Value.GetVideoFile()
		if !ok {
			return &api.PlayerStartPostBadRequest{}, nil
		}

		if file.Root >= len(config.Client.Media.Directories) {
			return &api.PlayerStartPostNotFound{}, nil
		}

		root := config.Client.Media.Directories[file.Root]
		openPath = path.Clean(path.Join(root, file.RelativePath))
		if !strings.HasPrefix(openPath, root) {
			return &api.PlayerStartPostNotFound{}, nil
		}

		fs, err := os.Stat(openPath)
		if err != nil || fs.IsDir() {
			return &api.PlayerStartPostNotFound{}, nil
		}
	case api.VideoStreamPlayerStartPostReq:
		openPath = string(req.Value.VideoStream)
	default:
		return &api.PlayerStartPostBadRequest{}, nil
	}
	return &api.PlayerStartPostAccepted{}, r.p.LoadFile(openPath)
}

// FilesGet implements GET /files operation.
// List file roots.
//
// GET /files
func (r *restServer) FilesGet(ctx context.Context) ([]string, error) {
	return config.Client.Media.Directories, nil
}

// GetFilesRootRelpath implements get-files-root-relpath operation.
//
// List files under the given root.
//
// GET /files/{root}
func (r *restServer) GetFilesRootRelpath(_ context.Context, params api.GetFilesRootRelpathParams) (api.GetFilesRootRelpathRes, error) {
	if params.Root >= len(config.Client.Media.Directories) {
		return &api.GetFilesRootRelpathNotFound{}, nil
	}

	rel := "."
	if params.Relpath.IsSet() {
		rel = params.Relpath.Value
	}

	p := path.Join(config.Client.Media.Directories[params.Root], params.Relpath.Value)
	entries, err := os.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %w", err)
	}

	var files []api.File
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue // ignore hidden files
		}
		files = append(files, api.File{
			Name:      entry.Name(),
			Directory: entry.IsDir(),
		})
	}

	return &api.Directory{
		Root:         params.Root,
		RelativePath: rel,
		Content:      files,
	}, nil
}

// GetPlayerPause implements get-player-pause operation.
// Query the player state on whether it is currently paused and provides the
// playback position if it
// is.
//
// GET /player/pause
func (r *restServer) GetPlayerPause(_ context.Context) (*api.Pause, error) {
	pause := r.p.GetPauseState()
	res := &api.Pause{
		Paused: pause,
	}
	if !pause {
		return res, nil
	}

	pos, err := r.p.GetPlaybackPos()
	if err != nil {
		return nil, fmt.Errorf("failed to get playback position: %w", err)
	}

	res.PausedMinusAt = api.NewOptPlaybackPosition(api.PlaybackPosition(pos.Milliseconds()) / 1000)
	return res, nil
}

// GetPlayerPosition implements get-player-position operation.
//
// Query the player for its current playback position.
//
// GET /player/position
func (r *restServer) GetPlayerPosition(_ context.Context) (api.PlaybackPosition, error) {
	pos, err := r.p.GetPlaybackPos()
	if err != nil {
		return 0, fmt.Errorf("failed to get playback position: %w", err)
	}

	return api.PlaybackPosition(pos.Milliseconds()) / 1000, nil
}

// GetStatus implements get-status operation.
//
// Returns information on the client currently used.
//
// GET /status
func (r *restServer) GetStatus(_ context.Context) (api.GetStatusRes, error) {
	re := regexp.MustCompile(`(?m)^.*?(\d)\.(\d)\.(\d)`)
	ver := buildinfo.VersionString()

	version := re.FindString(ver)
	vers := strings.SplitN(version, ".", 3)
	for i := len(vers); i < 3; i++ {
		vers = append(vers, "-1")
	}

	var versNumbers []int
	for _, v := range vers {
		part, _ := strconv.Atoi(v)
		versNumbers = append(versNumbers, part)
	}

	suffix := strings.TrimPrefix(buildinfo.VersionString(), version)
	status := *r.status

	switch status {
	case api.StatusConnectionOk:
		return &api.GetStatusOK{
			Connection: status,
			Player:     r.p.Name(),
			Version: api.Version{
				Major: versNumbers[0],
				Minor: versNumbers[1],
				Patch: versNumbers[2],
				Suffix: api.OptString{
					Value: suffix,
					Set:   len(suffix) > 0,
				},
			},
		}, nil
	default:
		return &api.GetStatusServiceUnavailable{
			Connection: status,
			Player:     r.p.Name(),
			Version: api.Version{
				Major: versNumbers[0],
				Minor: versNumbers[1],
				Patch: versNumbers[2],
				Suffix: api.OptString{
					Value: suffix,
					Set:   len(suffix) > 0,
				},
			},
		}, nil
	}
}

// PutPlayerPause implements put-player-pause operation.
//
// Set the player into the given pause state.
//
// PUT /player/pause
func (r *restServer) PutPlayerPause(_ context.Context, req bool) error {
	slog.Debug("setting pause state", "state", req)
	return r.p.Pause(req)
}

// PutPlayerPosition implements put-player-position operation.
//
// Seek to the specified position.
//
// PUT /player/position
func (r *restServer) PutPlayerPosition(_ context.Context, req api.PlaybackPosition) error {
	slog.Debug("seek", "seek_to", time.Duration(req)*time.Millisecond)
	return r.p.Seek(time.Duration(req * api.PlaybackPosition(time.Second)))
}
