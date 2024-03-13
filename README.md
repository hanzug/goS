## repository

### mysql

用户表

```
type User struct {
    UserID         int64  `gorm:"primarykey"`
    UserName       string `gorm:"unique"`
    NickName       string
    PasswordDigest string
}
```

输入数据表

```
type InputData struct {
    Id      int64  `gorm:"primarykey"`
    DocId   int64  `gorm:"index"`
    Title   string `gorm:"type:longtext"`
    Body    string `gorm:"type:longtext"`
    Url     string
    Score   float64
    Source  int
    IsIndex bool
}
```

收藏表

```
type Favorite struct {
    FavoriteID     int64             `gorm:"primarykey"` // 收藏夹id
    UserID         int64             `gorm:"index"`      // 用户id
    FavoriteName   string            `gorm:"unique"`     // 收藏夹名字
    FavoriteDetail []*FavoriteDetail `gorm:"many2many:f_to_fd;"`
}
```

收藏详细表

```
type FavoriteDetail struct {
    FavoriteDetailID int64       `gorm:"primarykey"`
    UserID           int64       // 用户id
    UrlID            int64       // url的id
    Url              string      // url地址
    Desc             string      // url的描述
    Favorite         []*Favorite `gorm:"many2many:f_to_fd;"`
}
```