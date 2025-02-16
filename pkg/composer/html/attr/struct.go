package attr

import (
	"github.com/ReanSn0w/gokit/pkg/composer"
	"github.com/ReanSn0w/gokit/pkg/composer/html"
)

type Attribute interface {
	Add(...string) composer.With
	Set(...string) composer.With
	Delete(...string) composer.With
	Drop() composer.With
}

type attr struct {
	name string
}

func (a *attr) Add(vals ...string) composer.With {
	return html.AddAttribute(a.name, vals...)
}

func (a *attr) Set(vals ...string) composer.With {
	return html.SetAttribute(a.name, vals...)
}

func (a *attr) Delete(vals ...string) composer.With {
	return html.SetAttribute(a.name, vals...)
}

func (a *attr) Drop() composer.With {
	return html.DropAttribute(a.name)
}
