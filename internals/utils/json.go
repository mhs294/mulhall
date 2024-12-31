package utils

import (
	"github.com/gin-gonic/gin"
)

// FromRequestJSON reads the JSON body of an HTTP request and deserializes its contents into the object
// referenced by the provided object pointer.
//
// v is the pointer to the object into which the JSON contents should be deserialized.
//
// ctx is the pointer to the [gin.Context] containing the HTTP request.
func FromRequestJSON(v interface{}, ctx *gin.Context) error {
	if err := ctx.ShouldBind(&v); err != nil {
		return err
	}

	return nil
}
