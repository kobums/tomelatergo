package user

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnEmail
    ColumnNickname
    ColumnSocialtype
    ColumnSocialid
    ColumnCreatedat
    ColumnUpdatedat
)

type Params struct {
    Column Column
    Value interface{}
}




