package operation

import "net/http"

func BadRequest(w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"reason": "Неверный формат запроса или его параметры"}`))
	return w
}

func InternalServerError(w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"reason": "Ошибка сервера"}`))
	return w
}

func Unauthorized(w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"reason": "Пользователь не существует или некорректен"}`))
	return w
}

func Forbidden(w http.ResponseWriter) http.ResponseWriter {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"reason": "Недостаточно прав для совершения действия"}`))
	return w
}

func WriteResponse(w http.ResponseWriter, code int, response []byte) http.ResponseWriter {
	w.WriteHeader(code)
	w.Write(response)
	return w
}
