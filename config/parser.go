package config

import (
	"fmt"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/vegacapsule/types"
	"github.com/hashicorp/go-cty-funcs/crypto"
	"github.com/hashicorp/go-cty-funcs/encoding"
	"github.com/hashicorp/go-cty-funcs/uuid"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/tryfunc"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var envFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "env_name",
			Type:             cty.String,
			AllowUnknown:     true,
			AllowDynamicType: false,
			AllowNull:        false,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		envName := args[0]
		if len(envName.AsString()) < 1 {
			return cty.NilVal, fmt.Errorf("the environment variable name cannot be empty")
		}

		return cty.StringVal(os.Getenv(envName.AsString())), nil
	},
})

func newEvalContext(genServices cty.Value) *hcl.EvalContext {
	return &hcl.EvalContext{
		Variables: map[string]cty.Value{"generated": genServices},
		Functions: map[string]function.Function{
			"abs":             stdlib.AbsoluteFunc,
			"base64decode":    encoding.Base64DecodeFunc,
			"base64encode":    encoding.Base64EncodeFunc,
			"bcrypt":          crypto.BcryptFunc,
			"ceil":            stdlib.CeilFunc,
			"chomp":           stdlib.ChompFunc,
			"chunklist":       stdlib.ChunklistFunc,
			"coalesce":        stdlib.CoalesceFunc,
			"coalescelist":    stdlib.CoalesceListFunc,
			"compact":         stdlib.CompactFunc,
			"concat":          stdlib.ConcatFunc,
			"contains":        stdlib.ContainsFunc,
			"csvdecode":       stdlib.CSVDecodeFunc,
			"distinct":        stdlib.DistinctFunc,
			"element":         stdlib.ElementFunc,
			"flatten":         stdlib.FlattenFunc,
			"floor":           stdlib.FloorFunc,
			"format":          stdlib.FormatFunc,
			"formatdate":      stdlib.FormatDateFunc,
			"formatlist":      stdlib.FormatListFunc,
			"indent":          stdlib.IndentFunc,
			"index":           stdlib.IndexFunc,
			"join":            stdlib.JoinFunc,
			"jsondecode":      stdlib.JSONDecodeFunc,
			"jsonencode":      stdlib.JSONEncodeFunc,
			"keys":            stdlib.KeysFunc,
			"length":          stdlib.LengthFunc,
			"log":             stdlib.LogFunc,
			"lookup":          stdlib.LookupFunc,
			"lower":           stdlib.LowerFunc,
			"max":             stdlib.MaxFunc,
			"md5":             crypto.Md5Func,
			"merge":           stdlib.MergeFunc,
			"min":             stdlib.MinFunc,
			"parseint":        stdlib.ParseIntFunc,
			"pow":             stdlib.PowFunc,
			"range":           stdlib.RangeFunc,
			"reverse":         stdlib.ReverseFunc,
			"replace":         stdlib.ReplaceFunc,
			"regex_replace":   stdlib.RegexReplaceFunc,
			"rsadecrypt":      crypto.RsaDecryptFunc,
			"setintersection": stdlib.SetIntersectionFunc,
			"setproduct":      stdlib.SetProductFunc,
			"setunion":        stdlib.SetUnionFunc,
			"sha1":            crypto.Sha1Func,
			"sha256":          crypto.Sha256Func,
			"sha512":          crypto.Sha512Func,
			"signum":          stdlib.SignumFunc,
			"slice":           stdlib.SliceFunc,
			"sort":            stdlib.SortFunc,
			"split":           stdlib.SplitFunc,
			"strlen":          stdlib.StrlenFunc,
			"strrev":          stdlib.ReverseFunc,
			"substr":          stdlib.SubstrFunc,
			"timeadd":         stdlib.TimeAddFunc,
			"title":           stdlib.TitleFunc,
			"trim":            stdlib.TrimFunc,
			"trimprefix":      stdlib.TrimPrefixFunc,
			"trimspace":       stdlib.TrimSpaceFunc,
			"trimsuffix":      stdlib.TrimSuffixFunc,
			"try":             tryfunc.TryFunc,
			"upper":           stdlib.UpperFunc,
			"urlencode":       encoding.URLEncodeFunc,
			"uuidv4":          uuid.V4Func,
			"uuidv5":          uuid.V5Func,
			"values":          stdlib.ValuesFunc,
			"zipmap":          stdlib.ZipmapFunc,
			"env":             envFunc,
		},
	}
}

func ApplyConfigContext(conf *Config, genServices *types.GeneratedServices) (*Config, error) {
	genServicesCtyVal, err := genServices.ToCtyValue()
	if err != nil {
		return nil, fmt.Errorf("failed to convert GeneratedServices to cty value: %w", err)
	}

	f, err := ParseHCLFile(conf.FilePath)
	if err != nil {
		return nil, err
	}

	decodeDiags := gohcl.DecodeBody(f.Body, newEvalContext(*genServicesCtyVal), conf)
	if decodeDiags.HasErrors() {
		return nil, fmt.Errorf("failed to decode config: %w", decodeDiags)
	}

	dir, _ := filepath.Split(conf.configDir)
	if err := conf.Validate(dir); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return conf, nil
}

func ParseHCLFile(filePath string) (*hcl.File, error) {
	parser := hclparse.NewParser()
	f, parseDiags := parser.ParseHCLFile(filePath)
	if parseDiags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL config file: %w", parseDiags)
	}

	return f, nil
}

func ParseConfigFile(filePath, outputDir string, genServices types.GeneratedServices) (*Config, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, err
	}

	if outputDir != "" {
		config.OutputDir = &outputDir
	}

	config.FilePath = filePath

	genServicesCtyVal, err := genServices.ToCtyValue()
	if err != nil {
		return nil, err
	}

	f, err := ParseHCLFile(filePath)
	if err != nil {
		return nil, err
	}

	decodeDiags := gohcl.DecodeBody(f.Body, newEvalContext(*genServicesCtyVal), config)
	if decodeDiags.HasErrors() {
		return nil, fmt.Errorf("failed to decode config: %w", decodeDiags)
	}

	dir, _ := filepath.Split(filePath)
	if err := config.Validate(dir); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}
	return config, nil
}
