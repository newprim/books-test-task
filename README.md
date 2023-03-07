# books-test-task

Выполнение [тестового задания](https://pastebin.com/Bq1u901B).

## Текст задания

Требуется реализовать Rest API с двумя endpoint'ами, с применением пакета
net/http или gin. "Content-Type" и "Accept" - "application/json".

1. Первый endpoint выводит список книг, состоящий из полей id, title, author,
   publisher_year методом GET /books. Кол-во книг 3.
2. Второй endpoint удаляет книгу из списка методом DELETE по маршруту /books и
   телом запроса {"book_Id": 1}.
3. В качестве хранилища данных предлагается использовать оперативную память.
   Хранилище должно быть выделено в отдельный слой, например repository.
4. Лимитировать доступ к endpoint'ам (throttling) не более 10RPS суммарно (для
   простоты, предлагаем буферизованные каналы). Способ реалиации не важен
   middleware или ServeHTTP уровень.
5. Необходимо разделить бизнес-логику и слой приложения (контроллеры),
   бизнес-логика должна быть независима от реализации (interface).
6. *Написать тесты на контроллеры.
7. Результат реализации скинуть ссылкой на github.

## Комментарий по реализации

Так как в задании ничего не сказано о retry логике, логгировании и т.д. - не
парился на этот счёт. Например, логгер просто
взял [отсюда](https://github.com/evrone/go-clean-template/blob/master/pkg/logger/logger.go).
HTTP сервер, ещё что-то по мелочи оттуда же.

Роутером я обычно пользуюсь `chi`, но так как сказано использовать или
стандартный, или `gin`, взял стандартный. Однако стандартный не поддерживает
миддлвейеры по типу `mux.Use(middleware)`, потому использую их не очень красиво.

Как я понял, в задании указана именно _рекомендация_ реализовывать троттлинг 
через каналы. Решил обойтись без каналов.

И последнее: я понял и реализовал троттлинг запросов так: если запросов
приходит больше чем 10 в секунду, все кроме 10 откидываются. Так как это
тестовое, спросить как нужно интерпретировать заданее, особо не у кого.
Если представленная реализация не удовлетворяет, прошу дать точно описание 
троттлинга.
