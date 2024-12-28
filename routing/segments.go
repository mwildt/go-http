package routing

import "strings"

type Parameters map[string]string

type Segment struct {
	value string
}

func (seg Segment) IsParam() (is bool, name string) {
	if strings.HasPrefix(seg.value, "{") && strings.HasSuffix(seg.value, "}") {
		return true, seg.value[1 : len(seg.value)-1]
	} else {
		return false, seg.value
	}
}

func (seg Segment) IsWildcard() bool {
	return seg.value == "*"
}

func (seg Segment) IsGlobalWildcard() bool {
	return seg.value == "**"
}

func (seg Segment) Print(params map[string]string) string {
	if isParam, paramName := seg.IsParam(); isParam {
		return params[paramName]
	} else if seg.IsWildcard() || seg.IsGlobalWildcard() {
		return ""
	} else {
		return seg.value
	}
}

type Segments []Segment

func NewSegments(template string) (segments Segments) {
	for _, value := range strings.Split(template, "/") {
		segments = append(segments, Segment{value})
	}
	return segments
}

func (segments Segments) Compare(path UriPath) (match bool, matched UriPath, params Parameters) {
	return compare(segments, path)
}

func (segments Segments) String() string {
	res := make([]string, 0)
	for _, segment := range segments {
		res = append(res, segment.value)
	}
	return strings.Join(res, "/")
}

func (segments Segments) Print(params map[string]string) string {
	uriPath := make(UriPath, 0)
	for _, segment := range segments {
		uriPath = append(uriPath, segment.Print(params))
	}
	return strings.Join(uriPath, "/")
}

func (segments Segments) Extend(path Segments) Segments {
	return NewSegments(segments.String() + path.String())
}

type UriPath []string

func NewUriPath(path string) UriPath {
	return strings.Split(path, "/")
}

func compare(segments Segments, path UriPath) (match bool, matched UriPath, params Parameters) {
	if len(segments) == 0 {
		return len(path) == 0, matched, make(Parameters)
	}

	if len(path) == 0 {
		return false, matched, make(Parameters)
	}

	if param, paramName := segments[0].IsParam(); param {
		if len(path[0]) == 0 {
			// a parameter always needs do have a non empty value
			return false, matched, make(Parameters)
		}

		match, matched, params = compare(segments[1:], path[1:])
		params[paramName] = path[0]
		return match, append(path[0:1], matched...), params

	} else if segments[0].IsWildcard() {
		match, matched, params = compare(segments[1:], path[1:])
		return match, append(path[0:1], matched...), params

	} else if segments[0].IsGlobalWildcard() {
		return true, path, make(Parameters)

	} else if segments[0].value == path[0] {
		match, matched, params = compare(segments[1:], path[1:])
		return match, append(path[0:1], matched...), params

	} else {
		// no match
		return false, matched, make(Parameters)
	}
}
