package utils

import (
	"crypto/md5"
	"fmt"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

var assetCache = make(map[string]string)

type CssUtils struct {
}

func NewCssUtils() *CssUtils {
	return &CssUtils{}
}

func (c *CssUtils) BuildCss() {
	result := api.Build(api.BuildOptions{
		EntryPoints:      []string{"./styles/index.css"},
		Outfile:          "./public/styles/main.css",
		Bundle:           true,
		Write:            true,
		LogLevel:         api.LogLevelInfo,
		MinifyWhitespace: true,
		MinifySyntax:     true,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}

func GetAssetHash(path string) string {
	if v, ok := assetCache[path]; ok {
		return v
	}

	data, err := os.ReadFile("public" + path[7:])
	if err != nil {
		return path 
	}

	hash := fmt.Sprintf("%x", md5.Sum(data))[:8]
	versionedPath := fmt.Sprintf("%s?v=%s", path, hash)

	assetCache[path] = versionedPath
	return versionedPath
}
