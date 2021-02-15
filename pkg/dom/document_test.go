package dom_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/djthorpe/data"
	"github.com/djthorpe/data/pkg/dom"
)

const (
	XMLFILE_A = "../../etc/xml/tiger.svg"
)

func Test_Document_001(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG, 0)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	if title := d.CreateElement("title"); title == nil {
		t.Error("Unexpected nil return from CreateElement")
	} else if err := d.AddChild(title); err != nil {
		t.Error(err)
	} else if err := title.AddChild(d.CreateText("Hello, World")); err != nil {
		t.Error(err)
	}

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\"><title>Hello, World</title></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_002(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG, 0)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	group := d.CreateElement("g")
	title := d.CreateElement("title")
	ta := d.CreateText("A")
	tb := d.CreateText("B")
	tc := d.CreateText("C")

	group.AddChild(title)
	group.AddChild(tc)
	d.AddChild(group)
	title.AddChild(ta)
	title.AddChild(tb)

	// Validate
	if len(d.Children()) != 1 || d.Children()[0] != group {
		t.Error("Unexpected children of document")
	}
	if len(group.Children()) != 2 || group.Children()[0] != title || group.Children()[1] != tc {
		t.Error("Unexpected children of group")
	}
	if len(title.Children()) != 2 || title.Children()[0] != ta || title.Children()[1] != tb {
		t.Error("Unexpected children of title")
	}

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\"><g><title>AB</title>C</g></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_003(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG, 0)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	group := d.CreateElement("group")
	ta := d.CreateText("A")
	tb := d.CreateText("B")
	tc := d.CreateText("C")
	d.AddChild(group)

	group.AddChild(ta)
	group.AddChild(tb)
	group.AddChild(tc)
	group.AddChild(ta) // order should be b,c,a

	// Validate
	if len(group.Children()) != 3 || group.Children()[0] != tb || group.Children()[1] != tc || group.Children()[2] != ta {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>BCA</group>" {
		t.Error("Unexpected: ", str)
	}

	group.AddChild(ta)
	group.AddChild(tb)
	group.AddChild(tc) // order should be a,b,c

	// Validate
	if len(group.Children()) != 3 || group.Children()[0] != ta || group.Children()[1] != tb || group.Children()[2] != tc {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>ABC</group>" {
		t.Error("Unexpected: ", str)
	}

	group.RemoveChild(ta) // order should be b,c

	// Validate
	if len(group.Children()) != 2 || group.Children()[0] != tb || group.Children()[1] != tc {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>BC</group>" {
		t.Error("Unexpected: ", str)
	}

	group.RemoveChild(ta)
	group.RemoveChild(tc) // order should be b

	// Validate
	if len(group.Children()) != 1 || group.Children()[0] != tb {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>B</group>" {
		t.Error("Unexpected: ", str)
	}

	group.AddChild(ta) // order should be b,a

	// Validate
	if len(group.Children()) != 2 || group.Children()[0] != tb || group.Children()[1] != ta {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>BA</group>" {
		t.Error("Unexpected: ", str)
	}

	group.AddChild(ta)
	group.AddChild(tb)
	group.AddChild(tc) // order should be a,b,c

	// Validate
	if len(group.Children()) != 3 || group.Children()[0] != ta || group.Children()[1] != tb || group.Children()[2] != tc {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>ABC</group>" {
		t.Error("Unexpected: ", str)
	}

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\"><group>ABC</group></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_004(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG, 0)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}

	// Add attribute
	if err := d.SetAttr("id", "document"); err != nil {
		t.Error(err)
	}

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"document\"></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_005(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG, 0)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}

	// Add attributes
	if err := d.SetAttr("id", "document"); err != nil {
		t.Error(err)
	}
	if err := d.SetAttr("id2", "document2"); err != nil {
		t.Error(err)
	}
	if err := d.SetAttr("id", "document3"); err != nil {
		t.Error(err)
	}

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"document3\" id2=\"document2\"></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_006(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG, 0)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}

	// Add image
	img := d.CreateElement("image")
	d.AddChild(img)

	if err := img.SetAttrNS("href", data.XmlNamespaceXLink, "http://www.google.com/"); err != nil {
		t.Fatal(err)
	}

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\"><image xmlns:xlink=\"http://www.w3.org/1999/xlink\" xlink:href=\"http://www.google.com/\"></image></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Get element by id
	if err := img.SetAttr("id", "test"); err != nil {
		t.Fatal(err)
	}
	if node := d.GetElementById("test"); node != img {
		t.Fatal("Unexpected return from GetElementById")
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_007(t *testing.T) {
	fh, err := os.Open(XMLFILE_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	// Read document
	if document, err := dom.Read(fh, data.DOMWriteDirective|data.DOMWriteIndentSpace2); err != nil {
		t.Fatal(err)
	} else {
		t.Log("\n" + fmt.Sprint(document))
	}
}

func Test_Document_008(t *testing.T) {
	d := dom.NewDocument("doc", 0)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	if comment := d.CreateComment("comment"); comment == nil {
		t.Error("Unexpected nil return from CreateElement")
	} else if err := d.AddChild(comment); err != nil {
		t.Error(err)
	}

	// Check output
	if str := fmt.Sprint(d); str != "<doc><!--comment--></doc>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_009(t *testing.T) {
	d := dom.NewDocument("doc", 0)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	if comment := d.CreateComment("comment"); comment == nil {
		t.Error("Unexpected nil return from CreateElement")
	} else if err := d.AddChild(comment); err != nil {
		t.Error(err)
	}

	// Check output
	if str := fmt.Sprint(d); str != "<doc><!--comment--></doc>" {
		t.Error("Unexpected: ", str)
	} else if _, err := dom.ReadEx(strings.NewReader(str), 0, func(node data.Node) error {
		t.Log("Validate=", node)
		return nil
	}); err != nil {
		t.Error("Validation Error: ", err)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}
