package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestEmailVerification(t *testing.T) {
	remoteWebDriver := "http://localhost:4444/wd/hub"

	wd, err := selenium.NewRemote(
		selenium.Capabilities{"browserName": "chrome"},
		remoteWebDriver,
	)
	if err != nil {
		t.Fatalf("[Ошибка] Не удалось создать веб-драйвер: %v", err)
	}
	defer wd.Quit()

	t.Log("[INFO] Веб-драйвер успешно создан")

	// Переход на страницу верификации
	t.Log("[INFO] Открываем страницу верификации email")
	err = wd.Get("http://localhost:8080/verify?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImRhbmEuanV6NDBAZ21haWwuY29tIiwiZXhwIjoxNzM4MzkwOTI0fQ.P1Ck3Mu0wz1YnFwOU_XHRJj6K96xJPIPysG-o1JtLxQ&email=dana.juz40@gmail.com")
	if err != nil {
		t.Fatalf("[Ошибка] Не удалось открыть страницу верификации: %v", err)
	}

	// Ожидаем редирект на страницу логина
	var currentURL string
	t.Log("[INFO] Ожидаем редирект на страницу логина")
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		var err error
		currentURL, err = wd.CurrentURL()
		if err != nil {
			return false, err
		}
		return strings.Contains(currentURL, "/login"), nil
	}, 10*time.Second)

	if err != nil {
		t.Fatalf("[Ошибка] Редирект не удался. Текущий URL: %v", currentURL)
	}

	t.Logf("[УСПЕХ] Email верифицирован, успешно редиректировано на: %s", currentURL)
}
