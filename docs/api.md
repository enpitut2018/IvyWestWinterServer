# User認証

## signup
```
curl --request POST \
  --url http://localhost:8080/signup \
  --data '{"userid": "ivy", "password": "password"}'
```

```js
{
    "message": "Success to create new user."
}
```

## signin
```
curl --request POST \
  --url http://localhost:8080/signin \
  --data '{"userid": "ivy", "password": "password"}'
```

```json
{
    "userid": "ivy",
    "abaterurl": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/user-face-photos/bfudnortqtp1qjpj7b9g.jpg",
    "token": "a735c3e8bc21cbe0f03e501a1529e0b4"
}
```

## User顔登録
```
curl --request POST \
  --url http://localhost:8080/uploadUserFace \
  --header 'Authorization: a735c3e8bc21cbe0f03e501a1529e0b4' \
  --data '{"source": "BASE64STRING"}'
```

```js
{
    "userid": "ivy",
    "abaterurl": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/user-face-photos/bfudnortqtp1qjpj7b9g.jpg",
    "token": "a735c3e8bc21cbe0f03e501a1529e0b4"
}
```

# 写真アップロード
```
curl --request POST \
  --url http://localhost:8080/uploads \
  --header 'Authorization: a735c3e8bc21cbe0f03e501a1529e0b4' \
  --data '{"source": "BASE64STRING"}'
```

```js
{
    "userid": "ivy",
    "url": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/upload-photos/bfudhurtqtp1rpmqpct0.jpg",
    "downloadUserIds": [
        "ivy",
        "sekine",
        "hiroto"
    ]
}
```

# 写真ダウンロード
```
curl --request GET \
  --url http://localhost:8080/downloads \
  --header 'Authorization: a735c3e8bc21cbe0f03e501a1529e0b4'
```

```js
[
    {
        "userid": "ivy",
        "url": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/upload-photos/bfudhurtqtp1rpmqpct0.jpg"
    }
]
```
