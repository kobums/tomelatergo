package question

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnContent
    ColumnMonth
    ColumnDay
    ColumnCreatedat
    ColumnUpdatedat
)

type Params struct {
    Column Column
    Value interface{}
}




