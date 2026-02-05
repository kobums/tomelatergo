package answer

type Column int

const (
    _ Column = iota
    
    ColumnId
    ColumnUid
    ColumnQid
    ColumnContent
    ColumnIspublic
    ColumnCreatedat
    ColumnUpdatedat
)

type Params struct {
    Column Column
    Value interface{}
}


type Ispublic int

const (
    _ Ispublic  = iota

    IspublicPrivate
    IspublicPublic
)

var Ispublics = []string{ "", "비공개", "공개" }



func GetIspublic(value Ispublic) string {
    i := int(value)
    if i <= 0 || i >= len(Ispublics) {
        return ""
    }
     
    return Ispublics[i]
}

func FindIspublic(value string) Ispublic {
    for i := 1; i < len(Ispublics); i++ {
        if Ispublics[i] == value {
            return Ispublic(i)
        }
    }
     
    return 0
}

func ConvertIspublic(value []int) []Ispublic {
     items := make([]Ispublic, 0)

     for item := range value {
         items = append(items, Ispublic(item))
     }
     
     return items
}

