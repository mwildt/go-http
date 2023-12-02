package http_utils

import (
	"testing"
)

func TestMatches(t *testing.T) {

	match, _, _ := compare(NewSegments("/api/context"), NewUriPath("/api/context"))
	Assert(t, match, "no match")

	match, _, _ = compare(NewSegments("/"), NewUriPath("/"))
	Assert(t, match, "no match")

	match, _, _ = compare(NewSegments("/*/*"), NewUriPath("/api/context"))
	Assert(t, match, "no match")

	match, _, _ = compare(NewSegments("/*/context"), NewUriPath("/api/context"))
	Assert(t, match, "no match")

	match, _, _ = compare(NewSegments("/api/context/"), NewUriPath("/api/context/"))
	Assert(t, match, "no match")

	match, _, _ = compare(NewSegments("/api/context/{id}"), NewUriPath("/api/context/1231"))
	Assert(t, match, "no match")

	match, _, _ = compare(NewSegments("/api/context/{id}/sub"), NewUriPath("/api/context/1231/sub"))
	Assert(t, match, "no match")

	match, _, _ = compare(NewSegments("/*"), NewUriPath("/"))
	Assert(t, match, "no match")

	match, _, _ = compare(NewSegments("/**"), NewUriPath("/this/is/all/matched"))
	Assert(t, match, "no match")

	match, _, _ = compare(NewSegments("/prefix/**"), NewUriPath("/prefix/this/is/all/matched"))
	Assert(t, match, "no match")

}

func TestNoMatches(t *testing.T) {

	match, _, _ := compare(NewSegments("/api/context"), NewUriPath("/api/"))
	Assert(t, !match, "match where no wanted")

	match, _, _ = compare(NewSegments("/api/context"), NewUriPath("/api/context/suffix"))
	Assert(t, !match, "match where no wanted")

	match, _, _ = compare(NewSegments("/api/context"), NewUriPath("/api/context/"))
	Assert(t, !match, "match where no wanted")

	match, _, _ = compare(NewSegments("/*"), NewUriPath("/api/context/"))
	Assert(t, !match, "match where no wanted")

	// ein parameter muss einen wert haben? JA
	match, _, _ = compare(NewSegments("/api/{id}}"), NewUriPath("/api/"))
	Assert(t, !match, "match where no wanted")

	// ein parameter muss einen wert haben? JA
	match, _, _ = compare(NewSegments("/api/{id}}/suffix"), NewUriPath("/api//suffix"))
	Assert(t, !match, "match where no wanted")

	match, _, _ = compare(NewSegments("/*"), NewUriPath("/only/first/is/matched"))
	Assert(t, !match, "match where no wanted")

}

func TestParameterExtraction(t *testing.T) {

	_, _, params := compare(NewSegments("/api/context"), NewUriPath("/api/context"))
	Assert(t, len(params) == 0, "parmas error")

	_, _, params = compare(NewSegments("/api/{id}}"), NewUriPath("/api/123"))
	Assert(t, params["id"] == "123", "parmas error")

	_, _, params = compare(NewSegments("/api/{id}}"), NewUriPath("/api/"))
	Assert(t, params["id"] == "", "parmas error")

}
