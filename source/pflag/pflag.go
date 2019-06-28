package pflag

import (
	"errors"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/boxgo/config/source"
	"github.com/spf13/pflag"
)

type (
	pflagsrc struct {
		opts source.Options
	}
)

func (pfs *pflagsrc) String() string {
	return "pflag"
}

func (pfs *pflagsrc) Watch() (source.Watcher, error) {
	return source.NewNoopWatcher()
}

func (pfs *pflagsrc) Read() (*source.ChangeSet, error) {
	if !pflag.Parsed() {
		return nil, errors.New("flags not parsed")
	}

	var changes map[string]interface{}

	visitFn := func(f *pflag.Flag) {
		n := strings.ToLower(f.Name)
		keys := strings.FieldsFunc(n, split)
		reverse(keys)

		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				tmp[k] = f.Value
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}

		mergo.Map(&changes, tmp) // need to sort error handling
		return
	}

	unset, ok := pfs.opts.Context.Value(includeUnsetKey{}).(bool)
	if ok && unset {
		pflag.VisitAll(visitFn)
	} else {
		pflag.Visit(visitFn)
	}

	b, err := pfs.opts.Encoder.Encode(changes)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Format:    pfs.opts.Encoder.String(),
		Data:      b,
		Timestamp: time.Now(),
		Source:    pfs.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

// NewSource returns a config source for integrating parsed pflags.
// Hyphens are delimiters for nesting, and all keys are lowercased.
//
// Example:
//      dbhost := flag.String("database-host", "localhost", "the db host name")
//
//      {
//          "database": {
//              "host": "localhost"
//          }
//      }
func NewSource(opts ...source.Option) source.Source {
	return &pflagsrc{opts: source.NewOptions(opts...)}
}

func split(r rune) bool {
	return r == '-' || r == '_' || r == '.'
}

func reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}
