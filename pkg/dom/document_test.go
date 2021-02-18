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

func CheckError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Document_001(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	if title := d.CreateElement("title"); title == nil {
		t.Error("Unexpected nil return from CreateElement")
	} else {
		CheckError(t, d.AddChild(title))
		CheckError(t, title.AddChild(d.CreateText("Hello, World")))
	}

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\"><title>Hello, World</title></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_002(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	group := d.CreateElement("g")
	title := d.CreateElement("title")
	ta := d.CreateText("A")
	tb := d.CreateComment("B")
	tc := d.CreateText("C")

	CheckError(t, group.AddChild(title))
	CheckError(t, group.AddChild(tc))
	CheckError(t, d.AddChild(group))
	CheckError(t, title.AddChild(ta))
	CheckError(t, title.AddChild(tb))

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
	if ta.PrevSibling() != nil {
		t.Error("Unexpected PrevSibling return", ta.PrevSibling())
	}
	if ta.NextSibling() != tb {
		t.Error("Unexpected NextSibling return", ta.NextSibling())
	}
	if tb.PrevSibling() != ta {
		t.Error("Unexpected PrevSibling return", tb.PrevSibling())
	}
	if tb.NextSibling() != nil {
		t.Error("Unexpected NextSibling return", tb.NextSibling())
	}

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\"><g><title>A<!--B--></title>C</g></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_003(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	group := d.CreateElement("group")
	ta := d.CreateText("A")
	tb := d.CreateText("B")
	tc := d.CreateText("C")
	CheckError(t, d.AddChild(group))
	CheckError(t, group.AddChild(ta))
	CheckError(t, group.AddChild(tb))
	CheckError(t, group.AddChild(tc))
	CheckError(t, group.AddChild(ta)) // order should be b,c,a

	// Validate
	if len(group.Children()) != 3 || group.Children()[0] != tb || group.Children()[1] != tc || group.Children()[2] != ta {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>BCA</group>" {
		t.Error("Unexpected: ", str)
	}

	CheckError(t, group.AddChild(ta))
	CheckError(t, group.AddChild(tb))
	CheckError(t, group.AddChild(tc)) // order should be a,b,c

	// Validate
	if len(group.Children()) != 3 || group.Children()[0] != ta || group.Children()[1] != tb || group.Children()[2] != tc {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>ABC</group>" {
		t.Error("Unexpected: ", str)
	}

	if ta.PrevSibling() != nil {
		t.Error("Unexpected return from PrevSibling")
	}
	if tb.PrevSibling() != ta {
		t.Error("Unexpected return from PrevSibling")
	}
	if tc.PrevSibling() != tb {
		t.Error("Unexpected return from PrevSibling")
	}
	if ta.NextSibling() != tb {
		t.Error("Unexpected return from PrevSibling")
	}
	if tb.NextSibling() != tc {
		t.Error("Unexpected return from PrevSibling")
	}
	if tc.NextSibling() != nil {
		t.Error("Unexpected return from PrevSibling")
	}

	CheckError(t, group.RemoveChild(ta)) // order should be b,c

	if ta.PrevSibling() != nil {
		t.Error("Unexpected return from PrevSibling")
	}
	if tb.PrevSibling() != nil {
		t.Error("Unexpected return from PrevSibling")
	}
	if tc.PrevSibling() != tb {
		t.Error("Unexpected return from PrevSibling")
	}
	if ta.NextSibling() != nil {
		t.Error("Unexpected return from PrevSibling")
	}
	if tb.NextSibling() != tc {
		t.Error("Unexpected return from PrevSibling")
	}
	if tc.NextSibling() != nil {
		t.Error("Unexpected return from PrevSibling")
	}

	// Validate
	if len(group.Children()) != 2 || group.Children()[0] != tb || group.Children()[1] != tc {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>BC</group>" {
		t.Error("Unexpected: ", str)
	}

	CheckError(t, group.RemoveChild(tc)) // order should be b

	// Validate
	if len(group.Children()) != 1 || group.Children()[0] != tb {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>B</group>" {
		t.Error("Unexpected: ", str)
	}

	CheckError(t, group.AddChild(ta)) // order should be b,a

	// Validate
	if len(group.Children()) != 2 || group.Children()[0] != tb || group.Children()[1] != ta {
		t.Fatal("Unexpected children of group", group.Children())
	}

	// Check output
	if str := fmt.Sprint(group); str != "<group>BA</group>" {
		t.Error("Unexpected: ", str)
	}

	CheckError(t, group.AddChild(ta))
	CheckError(t, group.AddChild(tb))
	CheckError(t, group.AddChild(tc)) // order should be a,b,c

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
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}

	// Add attribute
	CheckError(t, d.SetAttr("id", "document"))

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"document\"></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_005(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}

	// Add attributes
	CheckError(t, d.SetAttr("id", "document"))
	CheckError(t, d.SetAttr("id2", "document2"))
	CheckError(t, d.SetAttr("id", "document3"))

	// Check output
	if str := fmt.Sprint(d); str != "<svg xmlns=\"http://www.w3.org/2000/svg\" id=\"document3\" id2=\"document2\"></svg>" {
		t.Error("Unexpected: ", str)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_006(t *testing.T) {
	d := dom.NewDocumentNS("svg", data.XmlNamespaceSVG)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}

	// Add image
	img := d.CreateElement("image")
	CheckError(t, d.AddChild(img))

	if err := img.SetAttrNS("href", data.XmlNamespaceXLink, "http://www.google.com/"); err != nil {
		t.Fatal(err)
	}

	// Check output
	comp := `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><image xlink:href="http://www.google.com/"></image></svg>`
	if str := fmt.Sprint(d); str != comp {
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
	if document, err := dom.Read(fh); err != nil {
		t.Fatal(err)
	} else {
		t.Log("\n" + fmt.Sprint(document))
	}
}

func Test_Document_008(t *testing.T) {
	d := dom.NewDocument("doc")
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
	d := dom.NewDocument("doc")
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
	} else if _, err := dom.ReadEx(strings.NewReader(str), func(node data.Node) error {
		t.Log("Validate=", node)
		return nil
	}); err != nil {
		t.Error("Validation Error: ", err)
	}

	// Output
	t.Log("\n" + fmt.Sprint(d))
}

func Test_Document_010(t *testing.T) {
	fh, err := os.Open(XMLFILE_A)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	// Read document
	document, err := dom.Read(fh)
	if err != nil {
		t.Fatal(err)
	}

	// Obtain title element
	if nodes := document.GetElementsByTagNameNS("title", data.XmlNamespaceSVG); len(nodes) != 1 {
		t.Error("Expected one title tag, got ", nodes)
	} else if nodes[0].Name().Local != "title" {
		t.Error("Unexpected tag", nodes[0])
	}

	// Get group tags
	if nodes := document.GetElementsByTagNameNS("g", data.XmlNamespaceSVG); len(nodes) == 0 {
		t.Error("Expected more than one g tag, got ", nodes)
	}
}

func Test_Document_011(t *testing.T) {
	d := dom.NewDocument("xml")
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}

	group := d.CreateElement("g")
	ta := d.CreateElement("A")
	tb := d.CreateElement("B")
	tc := d.CreateElement("C")

	if group.Parent() != nil {
		t.Error("Unexpected Parent: ", group.Parent())
	}
	if ta.Parent() != nil {
		t.Error("Unexpected Parent: ", ta.Parent())
	}
	if tb.Parent() != nil {
		t.Error("Unexpected Parent: ", tb.Parent())
	}
	if tc.Parent() != nil {
		t.Error("Unexpected Parent: ", tc.Parent())
	}

	CheckError(t, group.AddChild(tb))
	CheckError(t, group.InsertChildBefore(ta, tb))
	CheckError(t, group.InsertChildBefore(tc, nil))

	// Validate
	if len(group.Children()) != 3 || group.Children()[0] != ta || group.Children()[1] != tb || group.Children()[2] != tc {
		t.Fatal("Unexpected children of group", group.Children())
	}
	if node := group.FirstChild(); node != ta {
		t.Fatal("Unexpected firstChild", node)
	}
	if node := group.LastChild(); node != tc {
		t.Fatal("Unexpected lastChild", node)
	}

	// Check output
	if str := fmt.Sprint(group); str != "<g><A></A><B></B><C></C></g>" {
		t.Error("Unexpected: ", str)
	}

	CheckError(t, group.AddChild(ta))
	CheckError(t, group.AddChild(tb))
	CheckError(t, group.AddChild(tc)) // order should be a,b,c

	// Validate
	if len(group.Children()) != 3 || group.Children()[0] != ta || group.Children()[1] != tb || group.Children()[2] != tc {
		t.Fatal("Unexpected children of group", group.Children())
	}
	if node := group.FirstChild(); node != ta {
		t.Fatal("Unexpected firstChild", node)
	}
	if node := group.LastChild(); node != tc {
		t.Fatal("Unexpected lastChild", node)
	}

	// Check output
	if str := fmt.Sprint(group); str != "<g><A></A><B></B><C></C></g>" {
		t.Error("Unexpected: ", str)
	}

	CheckError(t, group.RemoveAllChildren())

	// Validate
	if len(group.Children()) != 0 {
		t.Fatal("Unexpected children of group", group.Children())
	}
	if node := group.FirstChild(); node != nil {
		t.Fatal("Unexpected firstChild", node)
	}
	if node := group.LastChild(); node != nil {
		t.Fatal("Unexpected lastChild", node)
	}

	// Check output
	if str := fmt.Sprint(group); str != "<g></g>" {
		t.Error("Unexpected: ", str)
	}

	// Detached node
	if ta.Parent() != nil {
		t.Error("Unexpected Parent: ", ta.Parent())
	}
	if tb.Parent() != nil {
		t.Error("Unexpected Parent: ", tb.Parent())
	}
	if tc.Parent() != nil {
		t.Error("Unexpected Parent: ", tc.Parent())
	}
}

func Test_Document_012(t *testing.T) {
	d := dom.NewDocument("xml")
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}

	group := d.CreateElement("g")
	ta := d.CreateText("A")
	tb := d.CreateComment("B")
	tc := d.CreateElement("C")

	if group.Parent() != nil {
		t.Error("Unexpected Parent: ", group.Parent())
	}
	if ta.Parent() != nil {
		t.Error("Unexpected Parent: ", ta.Parent())
	}
	if tb.Parent() != nil {
		t.Error("Unexpected Parent: ", tb.Parent())
	}
	if tc.Parent() != nil {
		t.Error("Unexpected Parent: ", tc.Parent())
	}

	CheckError(t, group.AddChild(tb))
	CheckError(t, group.InsertChildBefore(ta, tb))
	CheckError(t, group.InsertChildBefore(tc, nil))

	// Validate
	if len(group.Children()) != 3 || group.Children()[0] != ta || group.Children()[1] != tb || group.Children()[2] != tc {
		t.Fatal("Unexpected children of group", group.Children())
	}
	if node := group.FirstChild(); node != ta {
		t.Fatal("Unexpected firstChild", node)
	}
	if node := group.LastChild(); node != tc {
		t.Fatal("Unexpected lastChild", node)
	}

	// Check output is ABC
	if str := fmt.Sprint(group); str != "<g>A<!--B--><C></C></g>" {
		t.Error("Unexpected: ", str)
	}

	CheckError(t, group.AddChild(ta))
	CheckError(t, group.AddChild(tb))
	CheckError(t, group.AddChild(tc)) // order should be a,b,c

	// Validate
	if len(group.Children()) != 3 || group.Children()[0] != ta || group.Children()[1] != tb || group.Children()[2] != tc {
		t.Fatal("Unexpected children of group", group.Children())
	}
	if node := group.FirstChild(); node != ta {
		t.Fatal("Unexpected firstChild", node)
	}
	if node := group.LastChild(); node != tc {
		t.Fatal("Unexpected lastChild", node)
	}

	// Check output
	if str := fmt.Sprint(group); str != "<g>A<!--B--><C></C></g>" {
		t.Error("Unexpected: ", str)
	}

	CheckError(t, group.RemoveAllChildren())

	// Validate
	if len(group.Children()) != 0 {
		t.Fatal("Unexpected children of group", group.Children())
	}
	if node := group.FirstChild(); node != nil {
		t.Fatal("Unexpected firstChild", node)
	}
	if node := group.LastChild(); node != nil {
		t.Fatal("Unexpected lastChild", node)
	}

	// Check output
	if str := fmt.Sprint(group); str != "<g></g>" {
		t.Error("Unexpected: ", str)
	}

	// Detached node
	if ta.Parent() != nil {
		t.Error("Unexpected Parent: ", ta.Parent())
	}
	if tb.Parent() != nil {
		t.Error("Unexpected Parent: ", tb.Parent())
	}
	if tc.Parent() != nil {
		t.Error("Unexpected Parent: ", tc.Parent())
	}
}

func Test_Document_013(t *testing.T) {
	d := dom.NewDocument("xml")
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}

	// Add attribute
	CheckError(t, d.SetAttr("version", "1.1"))
	if str := fmt.Sprint(d); str != "<xml version=\"1.1\"></xml>" {
		t.Error("Unexpected: ", str)
	}

	// Replace attribute
	CheckError(t, d.SetAttr("version", "1.2"))
	if str := fmt.Sprint(d); str != "<xml version=\"1.2\"></xml>" {
		t.Error("Unexpected: ", str)
	}

	// Add attribute
	CheckError(t, d.SetAttr("foo", "bar"))
	if str := fmt.Sprint(d); str != "<xml version=\"1.2\" foo=\"bar\"></xml>" {
		t.Error("Unexpected: ", str)
	}

	// Remove attribute
	CheckError(t, d.RemoveAttr("version"))
	if str := fmt.Sprint(d); str != "<xml foo=\"bar\"></xml>" {
		t.Error("Unexpected: ", str)
	}

	// Remove attribute
	CheckError(t, d.RemoveAttr("foo"))
	if str := fmt.Sprint(d); str != "<xml></xml>" {
		t.Error("Unexpected: ", str)
	}

	// Add attribute
	CheckError(t, d.SetAttr("foo2", "bar2"))
	if str := fmt.Sprint(d); str != "<xml foo2=\"bar2\"></xml>" {
		t.Error("Unexpected: ", str)
	}
}

func Test_Document_014(t *testing.T) {
	d := dom.NewDocumentNS("html", data.XmlNamespaceXHTML)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	if str := fmt.Sprint(d); str != `<html xmlns="http://www.w3.org/1999/xhtml"></html>` {
		t.Error("Unexpected return: ", str)
	} else {
		t.Log(str)
	}
}

func Test_Document_015(t *testing.T) {
	d := dom.NewDocumentNS("html", data.XmlNamespaceXHTML)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	CheckError(t, d.AddChild(d.CreateElementNS("head", data.XmlNamespaceXHTML)))
	if str := fmt.Sprint(d); str != `<html xmlns="http://www.w3.org/1999/xhtml"><head></head></html>` {
		t.Error("Unexpected return: ", str)
	} else {
		t.Log(str)
	}
}

func Test_Document_016(t *testing.T) {
	d := dom.NewDocumentNS("html", data.XmlNamespaceXHTML)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	CheckError(t, d.AddChild(d.CreateElementNS("head", data.XmlNamespaceSVG)))
	if str := fmt.Sprint(d); str != `<html xmlns="http://www.w3.org/1999/xhtml" xmlns:svg="http://www.w3.org/2000/svg"><svg:head></svg:head></html>` {
		t.Error("Unexpected return: ", str)
	} else {
		t.Log(str)
	}
}

func Test_Document_017(t *testing.T) {
	d := dom.NewDocumentNS("html", data.XmlNamespaceXHTML)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	CheckError(t, d.AddChild(d.CreateElementNS("head", data.XmlNamespaceSVG+" svg")))
	body := d.CreateElement("body")
	CheckError(t, body.SetAttr("id", "id1"))
	CheckError(t, body.SetAttrNS("id", data.XmlNamespaceSVG, "id2"))
	CheckError(t, d.AddChild(body))
	if str := fmt.Sprint(d); str != `<html xmlns="http://www.w3.org/1999/xhtml" xmlns:svg="http://www.w3.org/2000/svg"><svg:head></svg:head><body id="id1" svg:id="id2"></body></html>` {
		t.Error("Unexpected return: ", str)
	} else {
		t.Log(str)
	}
}

func Test_Document_018(t *testing.T) {
	d := dom.NewDocumentNS("html", data.XmlNamespaceXHTML)
	if d == nil {
		t.Fatal("Unexpected nil return from NewDocument")
	}
	CheckError(t, d.AddChild(d.CreateElementNS("head", data.XmlNamespaceSVG)))
	body := d.CreateElement("body")
	CheckError(t, body.SetAttr("id", "id1"))
	CheckError(t, body.SetAttrNS("id", data.XmlNamespaceSVG, "id2"))
	CheckError(t, d.AddChild(body))
	if str := fmt.Sprint(d); str != `<html xmlns="http://www.w3.org/1999/xhtml" xmlns:svg="http://www.w3.org/2000/svg"><svg:head></svg:head><body id="id1" svg:id="id2"></body></html>` {
		t.Error("Unexpected return: ", str)
	} else {
		t.Log(str)
	}
}

/*
	CheckError(t, d.AddChild(d.CreateElementNS("body", data.XmlNamespaceSVG+" svg")))
	CheckError(t, d.AddChild(d.CreateElementNS("div", data.XmlNamespaceXHTML)))
	CheckError(t, d.AddChild(d.CreateElementNS("div", data.XmlNamespaceSVG)))

	// I think the result should be something like:
	// <html xmlns="...">
	//   <head>...</head>
	//   <svg:body xmlns:svg="...">
	//     <div></div>
	//     <svg:div></svg:div>
	//   </svg:body>
	// </html>

	if str := fmt.Sprint(d); str != `<html xmlns="http://www.w3.org/1999/xhtml"><head></head></html>` {
		t.Error("Unexpected return: ", str)
	}
*/
