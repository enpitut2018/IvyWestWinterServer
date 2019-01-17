# IvyWestWinterServer

[FaceCafe](https://github.com/enpitut2018/IvyWestWinterFront)のAPIサーバです。

## エレベータピッチ

写真に写ってる人に共有するのが面倒なことと
自分の写っている写真を簡単に集めたいという問題を解決したい
自分の写真が欲しい人向けの
とった瞬間に自動で適切なユーザに共有される
写真共有サービスです。
これは顔認識でユーザの写真に写っている人を識別し、写っている人だけに自動で共有することができ、
GoogleフォトやLINEアルバムとは違って
圧倒的スピードで手間なく写真が共有できます。

## その他チーム用URL

* [プロダクトバックログ](https://github.com/enpitut2018/IvyWestWinterFront/projects/1)

* [フロントwiki](https://github.com/enpitut2018/IvyWestWinterFront/wiki)

* [サーバサイドwiki](https://github.com/enpitut2018/IvyWestWinterServer/wiki)

# How to use in local
```
>>> docker-compose build
>>> docker-compose up
```

# dockerのDBに入る。
```
>>> psql -h localhost -p 5432 -U postgres
```
