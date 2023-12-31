// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// FilesGet implements GET /files operation.
//
// List file roots.
//
// GET /files
func (UnimplementedHandler) FilesGet(ctx context.Context) (r []string, _ error) {
	return r, ht.ErrNotImplemented
}

// GetFilesRootRelpath implements get-files-root-relpath operation.
//
// List files under the given root.
//
// GET /files/{root}
func (UnimplementedHandler) GetFilesRootRelpath(ctx context.Context, params GetFilesRootRelpathParams) (r GetFilesRootRelpathRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetPlayerPause implements get-player-pause operation.
//
// Query the player state on whether it is currently paused and provides the playback position if it
// is.
//
// GET /player/pause
func (UnimplementedHandler) GetPlayerPause(ctx context.Context) (r *Pause, _ error) {
	return r, ht.ErrNotImplemented
}

// GetPlayerPosition implements get-player-position operation.
//
// Query the player for its current playback position.
//
// GET /player/position
func (UnimplementedHandler) GetPlayerPosition(ctx context.Context) (r PlaybackPosition, _ error) {
	return r, ht.ErrNotImplemented
}

// GetStatus implements get-status operation.
//
// Returns information on the client currently used.
//
// GET /status
func (UnimplementedHandler) GetStatus(ctx context.Context) (r GetStatusRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PlayerStartPost implements POST /player/start operation.
//
// Start playback of a video.
//
// POST /player/start
func (UnimplementedHandler) PlayerStartPost(ctx context.Context, req OptPlayerStartPostReq) (r PlayerStartPostRes, _ error) {
	return r, ht.ErrNotImplemented
}

// PutPlayerPause implements put-player-pause operation.
//
// Set the player into the given pause state.
//
// PUT /player/pause
func (UnimplementedHandler) PutPlayerPause(ctx context.Context, req bool) error {
	return ht.ErrNotImplemented
}

// PutPlayerPosition implements put-player-position operation.
//
// Seek to the specified position.
//
// PUT /player/position
func (UnimplementedHandler) PutPlayerPosition(ctx context.Context, req PlaybackPosition) error {
	return ht.ErrNotImplemented
}
