# potato
Небольшая копия редиса, написанная на го.

## Что есть
* Работа со строками/списками/хэшами и основные операции (keys, del). Все эти операции покрыты тестами.
* TTL для любого ключа (хотя само время удаления ключа может быть больше непосредственно заданного значения, ключ точно удалиться)
* Поддержка нескольких соединений
* Клиент и сервер общаются через сокеты - приложение клиента просто библиотечка-обёртка над протоколом

## Как запустить
> - cd potatoSlave && docker-compose up
> - cd potatoClient && go run main.go

## Как работает
* Четыре основные структуры: PotatoSlave, pstring, plist, pmap (последние три реализуют интерфейс potat).
* У potato slave есть мап в котором лежать potat'ы и мап с функциями, которые может вызывать клиент.
* TTL проверяется в отдельной горутине

## Как улучшать
Все улучшения, которые я вижу отмечены _TODO_ в коде. Из важного:
* Почти реализована авторизация, нужно только добавить логику проверки пароля в _authConnection_
* Можно сильно сократить число строк кода отрефакторив тесты и invocable функции (они однотипны)
* Сейчас _ttlCheckRoutine_ каждый раз проверяет все ключи на испорченность, кажется, что можно проверять каждый раз случайное подмножество, чтобы не иметь линейную по количеству ключей сложность.
