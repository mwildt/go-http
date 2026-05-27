package routing

import "strings"

// Parameters is a map of URL parameter names to their values.
type Parameters map[string]string

// Segment represents a single part of a URL path (e.g., "users", "{id}", "**").
type Segment struct {
	value string
}

// IsParam checks if the segment is a parameter (e.g., "{id}").
// Returns true and the parameter name if it is a parameter.
func (seg Segment) IsParam() (is bool, name string) {
	if strings.HasPrefix(seg.value, "{") && strings.HasSuffix(seg.value, "}") {
		return true, seg.value[1 : len(seg.value)-1]
	}
	return false, seg.value
}

// IsWildcard checks if the segment is a single wildcard ("*").
func (seg Segment) IsWildcard() bool {
	return seg.value == "*"
}

// IsGlobalWildcard checks if the segment is a global wildcard ("**").
func (seg Segment) IsGlobalWildcard() bool {
	return seg.value == "**"
}

// Print returns the string representation of the segment, replacing parameters with their values.
func (seg Segment) Print(params map[string]string) string {
	if isParam, paramName := seg.IsParam(); isParam {
		return params[paramName]
	} else if seg.IsWildcard() || seg.IsGlobalWildcard() {
		return ""
	}
	return seg.value
}

// Segments represents a URL path split into individual segments.
type Segments []Segment

// NewSegments creates a Segments slice from a URL path template.
// Example: NewSegments("/api/{id}") returns ["api", "{id}"].
// Empty segments are preserved to maintain consistency with NewUriPath().
func NewSegments(template string) (segments Segments) {
	for _, value := range strings.Split(template, "/") {
		segments = append(segments, Segment{value})
	}
	return segments
}

// Compare matches the segments against a UriPath, returning whether it matches,
// the matched path parts, and any extracted parameters.
func (segments Segments) Compare(path UriPath) (match bool, matched UriPath, params Parameters) {
	return compare(segments, path)
}

// String returns the segments joined as a URL path string.
func (segments Segments) String() string {
	res := make([]string, 0)
	for _, segment := range segments {
		res = append(res, segment.value)
	}
	return strings.Join(res, "/")
}

// Print returns the segments joined as a URL path string, with parameters replaced by their values.
func (segments Segments) Print(params map[string]string) string {
	uriPath := make(UriPath, 0)
	for _, segment := range segments {
		uriPath = append(uriPath, segment.Print(params))
	}
	return strings.Join(uriPath, "/")
}

func (segments Segments) Extend(path Segments) Segments {
	result := make(Segments, len(segments)+len(path))
	copy(result, segments)
	copy(result[len(segments):], path)
	return result
}

// UriPath represents a URL path split into individual parts.
type UriPath []string

// NewUriPath creates a UriPath from a URL path string.
// Empty segments are preserved to allow validation in compare().
func NewUriPath(path string) UriPath {
	parts := strings.Split(path, "/")
	result := make(UriPath, 0)
	for _, part := range parts {
		result = append(result, part)
	}
	return result
}

func compare(segments Segments, path UriPath) (match bool, matched UriPath, params Parameters) {
	params = make(Parameters)
	matched = make(UriPath, 0)
	i, j := 0, 0
	// Skip leading empty segments in both segments and path
	for i < len(segments) && len(segments[i].value) == 0 {
		i++
	}
	for j < len(path) && len(path[j]) == 0 {
		j++
	}
	for i < len(segments) && j < len(path) {
		// Reject empty segments in the middle of the path (e.g., from "//" in URL)
		if len(path[j]) == 0 {
			return false, matched, params
		}
		seg := segments[i]
		// Skip empty segments in segments (e.g., from "//" in template)
		if len(seg.value) == 0 {
			i++
			continue
		}
		if param, paramName := seg.IsParam(); param {
			if len(path[j]) == 0 {
				// A parameter must have a non-empty value
				return false, matched, params
			}
			params[paramName] = path[j]
			matched = append(matched, path[j])
			i++
			j++
		} else if seg.IsGlobalWildcard() {
			// Global wildcard matches the rest of the path (including empty)
			matched = append(matched, path[j:]...)
			return true, matched, params
		} else if seg.IsWildcard() {
			// Single wildcard matches one segment (including empty if at end)
			if j >= len(path) {
				// If no more path segments, wildcard can match empty
				matched = append(matched, "")
				i++
				j++
			} else {
				matched = append(matched, path[j])
				i++
				j++
			}
		} else if seg.value == path[j] {
			// Exact match
			matched = append(matched, path[j])
			i++
			j++
		} else {
			// No match
			return false, matched, params
		}
	}
	// Check if we consumed all segments and path parts (no skipping of trailing empty segments)
	return i == len(segments) && j == len(path), matched, params
}
