package jsonschema

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/qri-io/jsonschema"
	"github.com/tidwall/gjson"
)

// Errors
var (
	ErrReservedPath = errors.New("reserved path")
)

// Schema describes main interface for working with jsonschema validation.
type Schema interface {
	// Load fetches uri with jsonschema document which will be used for validating documents.
	Load(ctx context.Context, uri string) error
	// ValidateBytes validates byte slice with incoming json object.
	ValidateBytes(ctx context.Context, b []byte) ([]ValidationError, error)
}

type schema struct {
	schema        *jsonschema.Schema
	reservedPaths []string
}

func New(reservedPaths []string) Schema {
	return &schema{
		schema:        &jsonschema.Schema{},
		reservedPaths: reservedPaths,
	}
}

func (s *schema) checkReserved(bb []byte) error {
	if len(s.reservedPaths) == 0 {
		return nil
	}

	parsed := gjson.ParseBytes(bb)
	for _, path := range s.reservedPaths {
		if parsed.Get(path).Exists() {
			return errors.Wrap(ErrReservedPath, path)
		}
	}

	return nil
}

// Load fetches uri with jsonschema document which will be used for validating documents.
func (s *schema) Load(ctx context.Context, uri string) error {
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		err := jsonschema.FetchSchema(ctx, uri, s.schema)
		if err != nil {
			return errors.Wrapf(err, "fetch schema from remote uri: %s", uri)
		}
	}

	uri = strings.TrimPrefix(uri, "file://")

	u, err := url.Parse(uri)
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	body, err := ioutil.ReadFile(u.Path)
	if err != nil {
		return errors.Wrapf(err, "read file from path: %s", u.Path)
	}

	err = s.checkReserved(body)
	if err != nil {
		return errors.Wrap(err, "check presence of reserved paths")
	}

	if s.schema == nil {
		s.schema = &jsonschema.Schema{}
	}

	return json.Unmarshal(body, s.schema)
}

// ValidationError describes validation error
type ValidationError struct {
	// property that produced the error
	PropertyPath string
	// Message is a human-readable description of the error
	Message string
}

// ValidateBytes validates byte slice with incoming json object.
func (s *schema) ValidateBytes(ctx context.Context, b []byte) ([]ValidationError, error) {
	errs, err := s.schema.ValidateBytes(ctx, b)
	if err != nil {
		return nil, err
	}

	if len(errs) == 0 {
		return nil, nil
	}

	validErrs := make([]ValidationError, len(errs))
	for i := 0; i < len(errs); i++ {
		validErrs[i].PropertyPath = errs[i].PropertyPath
		validErrs[i].Message = errs[i].Message
	}

	return validErrs, err
}
