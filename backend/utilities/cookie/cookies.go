package cookie

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetCookie(c *gin.Context, name, value string, maxAge int, path, domain string, secure, httpOnly bool, sameSite http.SameSite) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:	 path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
	}
	http.SetCookie(c.Writer, &cookie)
	}