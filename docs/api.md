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

## User情報の一覧取得
```
curl --request GET \
  --url 'http://localhost:8080/users?ids=ivy,sekine,hiroto' \
  --header 'Authorization: a735c3e8bc21cbe0f03e501a1529e0b4' \
```

#### 備考
- idsで指定しなかった場合はアプリに登録している人全ての情報を取ってくる。
- idsを指定した場合はそのユーザーの情報のみ取ってくる。
  - idsに存在しないuseridを指定すると、そのユーザーの情報は何も返さない。

```
[
    {
        "userid": "hiroto",
        "avatarurl": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/user-face-photos/bfudnortqtp1qjpj7b9g.jpg"
    },
    {
        "userid": "sekine",
        "avatarurl": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/user-face-photos/bfudnortqtp1qjpj7b9g.jpg"
    },
    {
        "userid": "ivy",
        "avatarurl": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/user-face-photos/bfudnortqtp1qjpj7b9g.jpg"
    }
]
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

# 写真情報の取得

## 自分があげた写真の情報を返す
```
curl --request GET \
  --url 'http://localhost:8080/uploadPhotoInfos' \
  --header 'Authorization: a735c3e8bc21cbe0f03e501a1529e0b4' \
```

#### 備考
- useridsを指定する機能はない

```js
[
    {
        "id": 1,
        "url": "ivys-photo1.jpg",
        "uploader": {
            "id": "ivy",
            "avatar_url": ""
        },
        "downloaders": null
    },
    {
        "id": 2,
        "url": "ivys-photo2.jpg",
        "uploader": {
            "id": "ivy",
            "avatar_url": ""
        },
        "downloaders": [
            {
                "id": "ivy",
                "avatar_url": ""
            },
            {
                "id": "west",
                "avatar_url": ""
            }
        ]
    },
    {
        "id": 3,
        "url": "ivys-photo3.jpg",
        "uploader": {
            "id": "ivy",
            "avatar_url": ""
        },
        "downloaders": null
    },
    {
        "id": 4,
        "url": "ivys-photo4.jpg",
        "uploader": {
            "id": "ivy",
            "avatar_url": ""
        },
        "downloaders": null
    }
]
```


## 自分が写っている写真の情報を返す

```
curl --request GET \
  --url 'http://localhost:8080/downloadPhotoInfos?userids=west,sekine' \
  --header 'Authorization: a735c3e8bc21cbe0f03e501a1529e0b4' \
```

#### 備考
- 自分が写っている写真のみから選択する。
- useridを指定すると、自分とそのuseridが写っている写真のみ返すようになる。
- useridを指定しない場合は、自分が写っている写真を全て返す。
- 複数useridを指定した場合はor検索である。

``` js
[
    {
        "id": 2,
        "url": "ivys-photo2.jpg",
        "uploader": {
            "id": "ivy",
            "avatar_url": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/user-face-photos/bfudnortqtp1qjpj7b9g.jpg"
        },
        "downloaders": [
            {
                "id": "ivy",
                "avatar_url": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/user-face-photos/bfudnortqtp1qjpj7b9g.jpg"
            },
            {
                "id": "west",
                "avatar_url": "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/user-face-photos/bfudnortqtp1qjpj7b9g.jpg"
            }
        ]
    }
]
```
