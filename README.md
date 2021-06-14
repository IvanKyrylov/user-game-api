# Тестовое задание Go-developer:
Создать api для работы с данными пользователей с использованием базы данных
MongoDb.
Требуемый функционал:
Работа с пользователями:
- создать таблицу users и заполнить её тестовыми данными (ссылка ниже)
- разработать функционал запроса списка пользователей и информации о них
(статистика по сыгранным играм (сколько всего сыграно игр) + базовые данные)
- При выгрузке списка пользователей должна быть реализована постраничная
навигация
Работа с играми:
- создать таблицу user_games и заполнить её тестовыми данными для каждого
пользователя (тех что добавили в таблицу users с тестового набора), набор
тестовых данных по играм можно взять по ссылке ниже на каждого
пользователя выбирается из этого набора данных рандомное количество игр
(минимум 5000 игр на пользователя).
- также должна быть возможность получить статистику сгруппированную по
номерам игр и дням
- разработать функционал получения списка рейтинга пользователей (рейтинг
считается по количеству сыгранных пользователем игр), api должно отдавать данные
с постраничной навигацией.
Смысл тестового задания в том, чтобы увидеть навыки кандидата по работе со
структурами данных, оптимальном хранении данных и оптимизации запросов к БД.
Результаты выполнения разместить на github.
Само АПИ нужно задеплоить на Heroku.


## Выполнено 
Работа с пользователями:
- создать таблицу users и заполнить её тестовыми данными (ссылка ниже)
- разработать функционал запроса списка пользователей и информации о них
(статистика по сыгранным играм (сколько всего сыграно игр) + базовые данные)
- При выгрузке списка пользователей должна быть реализована постраничная
навигация
Работа с играми:
- создать таблицу user_games и заполнить её тестовыми данными для каждого
пользователя (тех что добавили в таблицу users с тестового набора), набор
тестовых данных по играм можно взять по ссылке ниже на каждого
пользователя выбирается из этого набора данных рандомное количество игр
(минимум 5000 игр на пользователя). Уточнение: В класторе MongoDB ограничение на 500мб, потому я сгенерировал 30-40 игр для каждого пользователя.
- также должна быть возможность получить статистику сгруппированную по
номерам игр и дням
- разработать функционал получения списка рейтинга пользователей (рейтинг
считается по количеству сыгранных пользователем игр), api должно отдавать данные
с постраничной навигацией.

Результаты выполнения разместить на github.
Само АПИ нужно задеплоить на Heroku.


## Описание 

### UserAPI
Получение списка пользователей (вся информация из таблицы) с ограничителем `{field limit}`- количество записей и `{page number}`- номер страницы:
`https://radiant-savannah-32971.herokuapp.com/api/users?limit={field limit}&page={page number}`

Получение даных о пользователе по id пользователя - `{UUID}` :
`https://radiant-savannah-32971.herokuapp.com/api/user/{UUID}`

Получение рейтинга пользователей (статистика по всем играм), вся информация про пользователя и значение рейтинга с ограничителем `{field limit}`- количество записей и `{page number}`- номер страницы:
`https://radiant-savannah-32971.herokuapp.com/api/users-rating?limit={field limit}&page={page number}`

### GameAPI
Получение списка игр (вся информация из таблицы) с ограничителем `{field limit}`- количество записей и `{page number}`- номер страницы:
`https://radiant-savannah-32971.herokuapp.com/api/games?limit={field limit}&page={page number}`

Получение даных о пользователе по id игры - `{ID}`:
`https://radiant-savannah-32971.herokuapp.com/api/user/{ID}`

Получение статистики сгруппированную по номерам игр и дням. `{UUID}`- id игрока, c`{start date}` - стартовая дата группировки, по `end date` - последняя дата группировки:
`https://radiant-savannah-32971.herokuapp.com/api/games-statistics?userId={UUID}&startDate={start date}&endDate={end date}`



## Url для теста(проверки)

### UserAPI

Список пользователей: https://radiant-savannah-32971.herokuapp.com/api/users?limit=10&page=0

Информация о пользователе: https://radiant-savannah-32971.herokuapp.com/api/user/60c60346a3807b4fba91fa09

Рейтинг пользователей: https://radiant-savannah-32971.herokuapp.com/api/users-rating?limit=1000&page=0


### GameAPI

Список игр: https://radiant-savannah-32971.herokuapp.com/api/games?limit=1000&page=0

Информация об играх: https://radiant-savannah-32971.herokuapp.com/api/game/60c6034ba3807b4fba9380a9

Статистики сгруппированную по номерам игр и дням: https://radiant-savannah-32971.herokuapp.com/api/games-statistics?userId=60c60346a3807b4fba91fa09&startDate=23-02-2019&endDate=04-03-2019


# Эпитафия
Извиняюсь что без Docker или make файлов, за ужасный код, отсутствие комментариев, отвратительную архитектуру, за readme на славянском. Но могу пока что, только так. Спасибо за внимание. Слава Омниссии!

