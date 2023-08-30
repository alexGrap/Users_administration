# avito backend intern
Проект, хранящий пользователя и тестируемые сегменты, в которых он состоит, а также историю добавления/удаления сегмента.

## Инструкции по запуску

```shell
git clone https://github.com/alexGrap/avito_intern.git
# клонирование репозитория

make run
# запуск сервиса и базы данных

http://localhost:3000/swagger/index.html#/
# документация к сервису
```

## Поддерживаемые запросы
``GET /getSegment`` - получение списка всех созданных сегментов</br>
``GET /getById?id={userId}`` - получение информации о сегментах, в которых состоит пользователь по его id.</br>
``POST /createSegment body{ "segmentName" : {segmentName}, "percent" : {usersPercent} }`` - создание нового сегмента с указанием процента пользователей, которые будут в него добавлены.</br>
``DELETE /deleteSegment body{ "segmentName" : {segmentName} }`` - удаление существующего сегмента.</br>
``PUT /subsription body{"userId" : {userId}, "add" : [ {adding segments} ], "delete" : [ {segment for delete} ]}`` - добавление и удаление ранее созданных сегментов пользователю</br>
``PUT /timeoutSubscribe body{"userId" : {userId}, "segmentName" : {segmentName}, "timeOut" : {timeout of living}}`` - добавление пользователю сегмента на срок {timeOut} дней (целочисленной значение).</br>
``GET /history?userId={userId}&from={date}&to={date}`` - возвращает информацию о добавлении и удалении пользователю сегментов в период с from до to дат.</br>

## Примеры запросов
```shell
curl -X GET localhost:3000/getSegment

curl -X GET localhost:3000/getById?id=1

curl -X POST localhost:3000/createSegment -H "Content-Type: application/json" -d '{ "segmentName" : "nkjwefkwe", "percents" : 50}'

curl -X DELETE localhost:3000/deleteSegment -H "Content-Type: application/json" -d '{ "segmentName" : "nkjwefkwe"}'

curl -X PUT localhost:3000/subscription -H "Content-Type: application/json" -d '{"userId" : 18, "add" : ["VOICE_MESSAGE"], "delete" : ["SUPER_SALE", "FIFTY_PERCENT_SALE"]}'

curl -X PUT localhost:3000/timeoutSubscribe -H "Content-Type: application/json" -d `{"userId" : 2, "segmentName" : "SUPER_SALE", "timeToDie" : 1}`

curl -X GET localhost:3000/history?userId=18&from=2020-03-20&to=2024-03-20
```