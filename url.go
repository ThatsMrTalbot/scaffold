package scaffold

import (
	"net/http"
	"path"
	"strings"

	"golang.org/x/net/context"
)

func pathSplit(p string) []string {
	p = strings.TrimSpace(p)
	if p == "" {
		return []string{}
	}

	p = path.Clean(p)
	parts := strings.Split(p, "/")
	if len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	return parts
}

// URLParts spliths a path into parts and caches it in the context
func URLParts(ctx context.Context, r *http.Request) (context.Context, []string) {
	if ctx == nil {
		ctx = context.Background()
	}

	if parts, ok := ctx.Value("scaffold_url_parts").([]string); ok {
		return ctx, parts
	}

	parts := pathSplit(r.URL.Path)

	ctx = context.WithValue(ctx, "scaffold_url_parts", parts)
	return URLParts(ctx, r)
}

// URLPart returns a part of the url and caches it in the context
func URLPart(ctx context.Context, r *http.Request, i int) (context.Context, string, bool) {
	ctx, parts := URLParts(ctx, r)
	if len(parts) > i && i >= 0 {
		return ctx, parts[i], true
	}
	return ctx, "", false
}
