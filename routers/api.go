package routers

import (
	"github.com/acaloiaro/frm"
)

type HttpRouter interface {
}

type Router interface {
	MountBuilder(router Router, f *frm.Frm)
}
