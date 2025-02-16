package attr

var (
	Class       Attribute = New("class")
	ID          Attribute = New("id")
	Href        Attribute = New("href")
	Src         Attribute = New("src")
	Alt         Attribute = New("alt")
	Title       Attribute = New("title")
	Style       Attribute = New("style")
	Name        Attribute = New("name")
	Value       Attribute = New("value")
	Type        Attribute = New("type")
	Placeholder Attribute = New("placeholder")
	Disabled    Attribute = New("disabled")
	ReadOnly    Attribute = New("readonly")
	Checked     Attribute = New("checked")
	Selected    Attribute = New("selected")
	Action      Attribute = New("action")
	Method      Attribute = New("method")
	Enctype     Attribute = New("enctype")
	Target      Attribute = New("target")
	Rel         Attribute = New("rel")
	Width       Attribute = New("width")
	Height      Attribute = New("height")
	Charset     Attribute = New("charset")
	Lang        Attribute = New("lang")
	TabIndex    Attribute = New("tabindex")
	AccessKey   Attribute = New("accesskey")
	Draggable   Attribute = New("draggable")
	Content     Attribute = New("content")
	HttpEquiv   Attribute = New("http-equiv")
)

func New(name string) Attribute {
	return &attr{name}
}
