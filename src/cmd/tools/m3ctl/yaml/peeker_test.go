// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package yaml

import (
	pb "github.com/m3db/m3/src/cmd/tools/m3ctl/yaml/generated"
	"io/ioutil"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
)

func TestPeekerPositive(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/basic_create.yaml")
	if err != nil {
		t.Fatalf("failed to read yaml test data:%v:\n", err)
	}
	urlpath, pbmessage, err := peeker(content)
	if err != nil {
		t.Fatalf("operation selector failed to encode the unknown operation yaml test data:%v:\n", err)
	}
	if urlpath != dbcreatePath {
		t.Errorf("urlpath is wrong:expected:%s:got:%s:\n", dbcreatePath, urlpath)
	}
	data, err := load(content, pbmessage)
	if err != nil {
		t.Fatalf("failed to encode to protocol:%v:\n", err)
	}
	dest := pb.DatabaseCreateRequestYaml{}
	unmarshaller := &jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err := unmarshaller.Unmarshal(data, &dest); err != nil {
		t.Fatalf("operation selector failed to unmarshal unknown operation data:%v:\n", err)
	}
	t.Logf("dest:%v:\n", dest)
	if dest.Request.NamespaceName != "default" {
		t.Errorf("dest NamespaceName does not have the correct value via operation:expected:%v:got:%v:", opCreate, dest.Request.NamespaceName)
	}
	if dest.Request.Type != "cluster" {
		t.Errorf("dest type does not have the correct value via operation:expected:%v:got:%v:", opCreate, dest.Request.Type)
	}
}
func TestPeekerNegative(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/unknown_operation.yaml")
	if err != nil {
		t.Fatalf("failed to read yaml test data:%v:\n", err)
	}
	if _, _, err := peeker(content); err == nil {
		t.Fatalf("operation selector should have returned an error\n")
	}
}
