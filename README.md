# IvyWestWinterServer
ivy-west winter server

# How to use
```
>>> docker-compose build
>>> docker-compose up
```

# dockerのDBに入る。
```
>>> psql -h localhost -p 5432 -U postgres
```

# URL
```
/signup
  - POST: ユーザーを作成する。
  {"userid": "ivy", "password": "pass"}

/signin
  - POST: Tokenを得る。
  {"userid": "ivy", "password": "pass"} -> token

/uploads(header: {"Autorization": token})
  - GET: ?userid="" ユーザがアップロードした画像をもらう。
  {"source": ""} -> {}
  - POST: 画像をアップロードする。
  {"source": ""} -> ok or fail
  - DELETE: アップロードした画像を消す。
  {"source": ""} -> ok or fail

/downloads(header: {"Autorization": token})
  - GET: ダウンロードした画像を得る。
  {"source": ""} -> {}
  - POST: 画像をダウンロードする。
  {"source": ""} -> ok or fail
  - DELETE: ダウンロードした画像を消す。
  {"source": ""} -> ok or fail
```
