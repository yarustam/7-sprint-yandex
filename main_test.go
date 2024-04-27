package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	//добавляю пакеты для тестирования
	"github.com/stretchr/testify/assert"  //позволяет проверить, выполняются ли условия, которые я задал
	"github.com/stretchr/testify/require" //пакет для тестирования, который останавливает выполнение, если условие не выполняется
)

// проверим что при запросе с количеством больше общего числа, возвращается ожидаемое число
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4 //ожидаемое количесвто запросов

	//создаю новый http-запрос методом get для пути /cafe
	req, err := http.NewRequest("GET", "/cafe", nil)
	if err != nil {
		require.NoError(t, err)
	}

	//получаю параметры запроса из url
	q := req.URL.Query()
	q.Add("count", "5")
	q.Add("city", "moscow")
	req.URL.RawQuery = q.Encode()

	//создаю переменную для сохранения ответа для тестирования http-обработчика
	responseRecorder := httptest.NewRecorder()

	//создаю обработчик http-запросов, который вызывает функцию mainHandle
	handler := http.HandlerFunc(mainHandle)

	//обрабатываю запрос с помощью обработчика и записываю результат в responseRecorder
	handler.ServeHTTP(responseRecorder, req)

	//получаю тело ответа в виде строки
	res := responseRecorder.Body.String()

	//разбиваю строку на массив, используя запятую в качестве разделителя
	arr := strings.Split(res, ",")

	//проверяю, что длинна массива = ожидаемому
	assert.Equal(t, totalCount, len(arr))
}

// эта функция отправляет запрос на эндпоинт /cafe с параметрами =4 и =omsk, а затем проверяет, что код ответа и тело соотвествуют ожидаемым значениям
func TestWrongCity(t *testing.T) {
	expectedError := "wrong city value"
	expectedCode := 400 //bab request ожидаемый код
	req, err := http.NewRequest("GET", "/cafe", nil)

	//проверяю есть ли ошибка
	if err != nil {
		require.NoError(t, err)
	}

	q := req.URL.Query()
	q.Add("count", "4")
	q.Add("city", "omsk")
	req.URL.RawQuery = q.Encode() // кодирую измененные параметры обратно в url запроса

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	resCode := responseRecorder.Code        //получаю фактический http-код ответа
	require.Equal(t, expectedCode, resCode) // проверяю, что фактический код = ожидаемому
	res := responseRecorder.Body.String()   //получаю тело ответа в виде строки
	require.Equal(t, res, expectedError)    //проверяю, что фактическое сообщение об ошибке соответсвует ожидаемому сообщению
}

func TestCorrectQuery(t *testing.T) {
	expectedCode := 200 //ожидаемый код
	req, err := http.NewRequest("GET", "/cafe", nil)
	if err != nil {
		require.NoError(t, err)
	}
	q := req.URL.Query()
	q.Add("count", "3")
	q.Add("city", "moscow")
	req.URL.RawQuery = q.Encode() //кодирую параметры запроса в URL

	responseRecorder := httptest.NewRecorder() //создаю запись ответа HTTP для захвата ответа от обработчика
	handler := http.HandlerFunc(mainHandle)    //определяю обработчик HTTP с именем mainHandle
	handler.ServeHTTP(responseRecorder, req)   //обрабатываю ответ и пишу его в responseRecorder
	resCode := responseRecorder.Code           //получаю фактический код состояния из записи ответа
	require.Equal(t, expectedCode, resCode)    //уточняю, что фактическое значение = ожидаемому, 200
	require.NotEmpty(t, responseRecorder.Body) //утверждаю, что тело ответа не пустое
}
