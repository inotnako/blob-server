# blob-server

для тестового задания понадобится
go - https://golang.org/dl/
mongodb - https://www.mongodb.org/downloads
git+ акк в github.com(чтобы туда потом выложить)
любой редактор (мы используем этот https://www.jetbrains.com/idea/download/ + плюс к нему плагин под го https://github.com/go-lang-plugin-org там есть фактически все что нужно для автокоплит, переход к определению, просмотр доки, запуск и тд)

# Задание

написать сервис по хранению файлов в mongodb.gridfs(http://docs.mongodb.org/manual/core/gridfs/)
то есть это http сервис со след api:

```
POST /api/v1/file 
GET /api/v1/file/:id
GET /api/v1/file - list files
DELETE /api/v1/file/:id - remove
```

