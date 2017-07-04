package main

import (
	"net/http"

	"github.com/herb-go/herb/model"
	"github.com/herb-go/herb/render"
)

func exampleApiAction(w http.ResponseWriter, r *http.Request) {
	form := exampleFormModel{}
	if model.MustValidateJSONPost(r, &form) {
		render.JSON(w, form.Errors(), 200)
	} else {
		model.MustRenderErrors(w, &form)
	}
}
