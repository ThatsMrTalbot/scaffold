package encoding

import (
	"sort"
	"strconv"
	"strings"
)

type mimeInfo struct {
	mime string
	typ  string
	sub  string
	q    float32
}

func (m *mimeInfo) equals(info *mimeInfo) bool {
	typ := info.typ == m.typ || info.typ == "*"
	sub := info.sub == m.sub || info.sub == "*"
	return typ && sub
}

type mimeInfoList []*mimeInfo

func (a mimeInfoList) Len() int {
	return len(a)
}

func (a mimeInfoList) Less(i, j int) bool {
	return a[i].q > a[j].q
}

func (a mimeInfoList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func parseMimeInfo(mime string) *mimeInfo {
	mime = strings.TrimSpace(mime)
	i := strings.Index(mime, "/")
	if i == -1 {
		return nil
	}

	return &mimeInfo{
		mime: mime,
		typ:  mime[:i],
		sub:  mime[i+1:],
		q:    1,
	}
}

func parseAcceptHeader(accept string) []*mimeInfo {
	accept = strings.TrimSpace(accept)
	if accept == "" {
		return nil
	}
	var mimes mimeInfoList

	for {
		buf := accept
		i := strings.Index(accept, ",")

		if i != -1 {
			accept, buf = accept[i+1:], accept[:i]
		}

		if a := parseAcceptMime(buf); a != nil {
			mimes = append(mimes, a)
		}

		if i == -1 {
			break
		}
	}

	sort.Sort(mimes)

	return []*mimeInfo(mimes)
}

func parseAcceptMime(a string) *mimeInfo {
	// Parse mime as a whole
	mime := ""
	buf := string(a)
	i := strings.Index(buf, ";")

	if i != -1 {
		mime, buf = strings.TrimSpace(buf[:i]), buf[i+1:]
	} else {
		mime = strings.TrimSpace(buf)
	}

	// Parse type and subtype
	info := parseMimeInfo(mime)
	if info == nil {
		return nil
	}

	// Parse q
	param := ""
	for {
		i := strings.Index(buf, ";")
		if i == -1 {
			param, buf = strings.TrimSpace(buf), ""
		} else {
			param, buf = strings.TrimSpace(buf[:i]), buf[i+1:]
		}

		if param[:2] == "q=" {
			f, _ := strconv.ParseFloat(param[2:], 32)
			info.q = float32(f)
			break
		}

		if buf == "" {
			break
		}
	}

	return info
}
