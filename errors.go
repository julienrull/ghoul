package ghoul

import (
	"fmt"
	"net/http"
)


func (e APIError) Error() string {
    return  fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error) APIError {
    return APIError{
        StatusCode: statusCode,
        Msg:        err.Error(),
    }
}

func InvalidRequestData(errors map[string]string) APIError {
    return APIError{
        StatusCode: http.StatusUnprocessableEntity,
        Msg:        errors,
    }
}

func InvalidJSON() APIError {
    return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid JSON request data"))
}

var DefaultErrorHandler ErrorHandler = func (ctx Ctx, err error) error {
    if err != nil {
        if apiErr, ok := err.(APIError); ok {
            ctx.Status(apiErr.StatusCode)
            ctx.Json(apiErr)
        }else {
            ctx.Status(http.StatusInternalServerError)
            errResp := map[string]any{
                "statusCode": http.StatusInternalServerError,
                "msg": http.StatusText(http.StatusInternalServerError),
            }
            ctx.Json(errResp)
        }
    }
    return err
}

var HandleError = DefaultErrorHandler


