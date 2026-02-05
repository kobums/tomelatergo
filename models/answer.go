package models

import (
    "tomelater/global/config"
    "tomelater/models/answer"
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "tomelater/global/log"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"

)

type Answer struct {
            
    Id                int64 `json:"id"`         
    Uid                int64 `json:"uid"`         
    Qid                int64 `json:"qid"`         
    Content                string `json:"content"`         
    Ispublic                answer.Ispublic `json:"ispublic"`         
    Createdat                string `json:"createdat"`         
    Updatedat                string `json:"updatedat"` 
    
    Extra                    map[string]interface{} `json:"extra"`
}

type AnswerManager struct {
    Conn    *Connection
    Result  *sql.Result
    Index   string
    Isolation   bool
    SelectQuery  string
    JoinQuery string
    CountQuery   string
    GroupQuery string
    SelectLog bool
    Log bool
}

func (c *Answer) AddExtra(key string, value interface{}) {    
	c.Extra[key] = value     
}

func NewAnswerManager(conn *Connection) *AnswerManager {
    var item AnswerManager


    if conn == nil {
        item.Conn = NewConnection()
        item.Isolation = false
    } else {
        item.Conn = conn 
        item.Isolation = conn.Isolation
    }

    item.Index = ""
    item.SelectLog = config.Log.Database
    item.Log = config.Log.Database

    return &item
}

func (p *AnswerManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *AnswerManager) SetIndex(index string) {
    p.Index = index
}

func (p *AnswerManager) SetCountQuery(query string) {
    p.CountQuery = query
}

func (p *AnswerManager) SetSelectQuery(query string) {
    p.SelectQuery = query
}

func (p *AnswerManager) Exec(query string, params ...interface{}) (sql.Result, error) {
    if p.Log {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Exec(query, params...)
}

func (p *AnswerManager) Query(query string, params ...interface{}) (*sql.Rows, error) {
    if p.Isolation {
        query += " for update"
    }

    if p.SelectLog {
       if len(params) > 0 {
	       log.Debug().Str("query", query).Any("param", params).Msg("SQL")
       } else {
	       log.Debug().Str("query", query).Msg("SQL")
       }
    }

    return p.Conn.Query(query, params...)
}

func (p *AnswerManager) GetQuery() string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder

    ret.WriteString("select a_id, a_uid, a_qid, a_content, a_ispublic, a_createdat, a_updatedat, u_id, u_email, u_nickname, u_socialtype, u_socialid, u_createdat, u_updatedat, q_id, q_content, q_month, q_day, q_createdat, q_updatedat from answer_tb, user_tb, question_tb")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    if p.JoinQuery != "" {
        ret.WriteString(", ")
        ret.WriteString(p.JoinQuery)
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and a_uid = u_id ")
    
    ret.WriteString("and a_qid = q_id ")
    

    return ret.String()
}

func (p *AnswerManager) GetQuerySelect() string {
    if p.CountQuery != "" {
        return p.CountQuery    
    }

    var ret strings.Builder
    
    ret.WriteString("select count(*) from answer_tb")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    if p.JoinQuery != "" {
        ret.WriteString(", ")
        ret.WriteString(p.JoinQuery)
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and a_uid = u_id ")
    
    ret.WriteString("and a_qid = q_id ")
    

    return ret.String()
}

func (p *AnswerManager) GetQueryGroup(name string) string {
    if p.SelectQuery != "" {
        return p.SelectQuery    
    }

    var ret strings.Builder
    ret.WriteString("select a_")
    ret.WriteString(name)
    ret.WriteString(", count(*) from answer_tb ")

    if p.Index != "" {
        ret.WriteString(" use index(")
        ret.WriteString(p.Index)
        ret.WriteString(")")
    }

    ret.WriteString(" where 1=1 ")
    
    ret.WriteString("and a_uid = u_id ")
    
    ret.WriteString("and a_qid = q_id ")
    

    return ret.String()
}

func (p *AnswerManager) Truncate() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    query := "truncate answer_tb "
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return nil
}

func (p *AnswerManager) Insert(item *Answer) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now().UTC().Add(time.Hour * 9)
        //t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    
    if item.Createdat == "" {
       item.Createdat = "1000-01-01 00:00:00"
    }
	
    if item.Updatedat == "" {
       item.Updatedat = "1000-01-01 00:00:00"
    }
	

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into answer_tb (a_id, a_uid, a_qid, a_content, a_ispublic, a_createdat, a_updatedat) values (?, ?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Id, item.Uid, item.Qid, item.Content, item.Ispublic, item.Createdat, item.Updatedat)
    } else {
        query = "insert into answer_tb (a_uid, a_qid, a_content, a_ispublic, a_createdat, a_updatedat) values (?, ?, ?, ?, ?, ?)"
        res, err = p.Exec(query, item.Uid, item.Qid, item.Content, item.Ispublic, item.Createdat, item.Updatedat)
    }
    
    if err == nil {
        p.Result = &res
        
    } else {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
        p.Result = nil
    }

    return err
}

func (p *AnswerManager) Delete(id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from answer_tb where a_id = ?"
    _, err := p.Exec(query, id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

func (p *AnswerManager) DeleteAll() error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query := "delete from answer_tb"
    _, err := p.Exec(query)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *AnswerManager) MakeQuery(initQuery string , postQuery string, initParams []interface{}, args []interface{}) (string, []interface{}) {
    var params []interface{}
    if initParams != nil {
        params = append(params, initParams...)
    }

    pos := 1

    var query strings.Builder
	query.WriteString(initQuery)

    for _, arg := range args {
        switch v := arg.(type) {        
        case Where:
            item := v

            if strings.Contains(item.Column, "_") {
                query.WriteString(" and ")
            } else {
                query.WriteString(" and a_")
            }
            query.WriteString(item.Column)

            if item.Compare == "in" {
                query.WriteString(" in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "not in" {
                query.WriteString(" not in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "between" {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(fmt.Sprintf(" between $%v and $%v", pos, pos + 1))
                    pos += 2
                } else {
                    query.WriteString(" between ? and ?")
                }

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(fmt.Sprintf(" $%v", pos))
                    pos++
                } else {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(" ?")
                }
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        case Custom:
             item := v

            query.WriteString(" and ")
            query.WriteString(item.Query)
        }        
    }

	query.WriteString(postQuery)

    return query.String(), params
}

func (p *AnswerManager) DeleteWhere(args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    query, params := p.MakeQuery("delete from answer_tb where 1=1", "", nil, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    return err
}

func (p *AnswerManager) Update(item *Answer) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }
    
    
    if item.Createdat == "" {
       item.Createdat = "1000-01-01 00:00:00"
    }
	
    if item.Updatedat == "" {
       item.Updatedat = "1000-01-01 00:00:00"
    }
	

	query := "update answer_tb set a_uid = ?, a_qid = ?, a_content = ?, a_ispublic = ?, a_createdat = ?, a_updatedat = ? where a_id = ?"
	_, err := p.Exec(query, item.Uid, item.Qid, item.Content, item.Ispublic, item.Createdat, item.Updatedat, item.Id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }
    
        
    return err
}

func (p *AnswerManager) UpdateWhere(columns []answer.Params, args []interface{}) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

    var initQuery strings.Builder
    var initParams []interface{}

    initQuery.WriteString("update answer_tb set ")
    for i, v := range columns {
        if i > 0 {
            initQuery.WriteString(", ")
        }

        if v.Column == answer.ColumnId {
        initQuery.WriteString("a_id = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == answer.ColumnUid {
        initQuery.WriteString("a_uid = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == answer.ColumnQid {
        initQuery.WriteString("a_qid = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == answer.ColumnContent {
        initQuery.WriteString("a_content = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == answer.ColumnIspublic {
        initQuery.WriteString("a_ispublic = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == answer.ColumnCreatedat {
        initQuery.WriteString("a_createdat = ?")
        initParams = append(initParams, v.Value)
        } else if v.Column == answer.ColumnUpdatedat {
        initQuery.WriteString("a_updatedat = ?")
        initParams = append(initParams, v.Value)
        } else {
        
        }
    }

    initQuery.WriteString(" where 1=1 ")

    query, params := p.MakeQuery(initQuery.String(), "", initParams, args)
    _, err := p.Exec(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
    }

    
    return err
}

/*


func (p *AnswerManager) UpdateUid(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update answer_tb set a_uid = ? where a_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AnswerManager) UpdateQid(value int64, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update answer_tb set a_qid = ? where a_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AnswerManager) UpdateContent(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update answer_tb set a_content = ? where a_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AnswerManager) UpdateIspublic(value answer.Ispublic, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update answer_tb set a_ispublic = ? where a_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AnswerManager) UpdateCreatedat(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update answer_tb set a_createdat = ? where a_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}

func (p *AnswerManager) UpdateUpdatedat(value string, id int64) error {
    if !p.Conn.IsConnect() {
        return errors.New("Connection Error")
    }

	query := "update answer_tb set a_updatedat = ? where a_id = ?"
	_, err := p.Exec(query, value, id)

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
    }

    return err
}


*/

func (p *AnswerManager) GetIdentity() int64 {
    if !p.Conn.IsConnect() {
        return 0
    }

    id, err := (*p.Result).LastInsertId()

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
        return 0
    } else {
        return id
    }
}

func (p *Answer) InitExtra() {
    p.Extra = map[string]interface{}{
            "ispublic":     answer.GetIspublic(p.Ispublic),

    }
}

func (p *AnswerManager) ReadRow(rows *sql.Rows) *Answer {
    var item Answer
    var err error

    var _user User
    var _question Question
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Uid, &item.Qid, &item.Content, &item.Ispublic, &item.Createdat, &item.Updatedat, &_user.Id, &_user.Email, &_user.Nickname, &_user.Socialtype, &_user.Socialid, &_user.Createdat, &_user.Updatedat, &_question.Id, &_question.Content, &_question.Month, &_question.Day, &_question.Createdat, &_question.Updatedat)
        
        if item.Createdat == "0000-00-00 00:00:00" || item.Createdat == "1000-01-01 00:00:00" || item.Createdat == "9999-01-01 00:00:00" {
            item.Createdat = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Createdat = strings.ReplaceAll(strings.ReplaceAll(item.Createdat, "T", " "), "Z", "")
        }
		
        if item.Updatedat == "0000-00-00 00:00:00" || item.Updatedat == "1000-01-01 00:00:00" || item.Updatedat == "9999-01-01 00:00:00" {
            item.Updatedat = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Updatedat = strings.ReplaceAll(strings.ReplaceAll(item.Updatedat, "T", " "), "Z", "")
        }
		

    } else {
        return nil
    }

    if err != nil {
        if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
        }
        return nil
    } else {

        item.InitExtra()
        _user.InitExtra()
        item.AddExtra("user",  _user)
_question.InitExtra()
        item.AddExtra("question",  _question)

        return &item
    }
}

func (p *AnswerManager) ReadRows(rows *sql.Rows) []Answer {
    var items []Answer

    for rows.Next() {
        var item Answer
        var _user User
        var _question Question
        

        err := rows.Scan(&item.Id, &item.Uid, &item.Qid, &item.Content, &item.Ispublic, &item.Createdat, &item.Updatedat, &_user.Id, &_user.Email, &_user.Nickname, &_user.Socialtype, &_user.Socialid, &_user.Createdat, &_user.Updatedat, &_question.Id, &_question.Content, &_question.Month, &_question.Day, &_question.Createdat, &_question.Updatedat)
        if err != nil {
           if p.Log {
             log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        
        if item.Createdat == "0000-00-00 00:00:00" || item.Createdat == "1000-01-01 00:00:00" || item.Createdat == "9999-01-01 00:00:00" {
            item.Createdat = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Createdat = strings.ReplaceAll(strings.ReplaceAll(item.Createdat, "T", " "), "Z", "")
        }
		
        if item.Updatedat == "0000-00-00 00:00:00" || item.Updatedat == "1000-01-01 00:00:00" || item.Updatedat == "9999-01-01 00:00:00" {
            item.Updatedat = ""
        }

        if config.Database.Type == config.Postgresql {
            item.Updatedat = strings.ReplaceAll(strings.ReplaceAll(item.Updatedat, "T", " "), "Z", "")
        }
		

        item.InitExtra()
        _user.InitExtra()
        item.AddExtra("user",  _user)
_question.InitExtra()
        item.AddExtra("question",  _question)

        items = append(items, item)
    }


     return items
}

func (p *AnswerManager) Get(id int64) *Answer {
    if !p.Conn.IsConnect() {
        return nil
    }

    var query strings.Builder
    query.WriteString(p.GetQuery())
    query.WriteString(" and a_id = ?")

    
    query.WriteString(" and a_uid = u_id")
    
    query.WriteString(" and a_qid = q_id")
    
    
    rows, err := p.Query(query.String(), id)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
       return nil
    }

    defer rows.Close()

    return p.ReadRow(rows)
}

func (p *AnswerManager) GetWhere(args []interface{}) *Answer {
    items := p.Find(args)
    if len(items) == 0 {
        return nil
    }

    return &items[0]
}

func (p *AnswerManager) Count(args []interface{}) int {
    if !p.Conn.IsConnect() {
        return 0
    }

    query, params := p.MakeQuery(p.GetQuerySelect(), p.GroupQuery, nil, args)
    rows, err := p.Query(query, params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
       return 0
    }

    defer rows.Close()

    if !rows.Next() {
        return 0
    }

    cnt := 0
    err = rows.Scan(&cnt)

    if err != nil {
        return 0
    } else {
        return cnt
    }
}

func (p *AnswerManager) FindAll() []Answer {
    return p.Find(nil)
}

func (p *AnswerManager) Find(args []interface{}) []Answer {
    if !p.Conn.IsConnect() {
        var items []Answer
        return items
    }

    var params []interface{}
    baseQuery := p.GetQuery()

    var query strings.Builder

    page := 0
    pagesize := 0
    orderby := ""

    pos := 1
    
    for _, arg := range args {
        switch v := arg.(type) {
        case PagingType:
            item := v
            page = item.Page
            pagesize = item.Pagesize            
        case OrderingType:
            item := v
            orderby = item.Order
        case LimitType:
            item := v
            page = 1
            pagesize = item.Limit
        case OptionType:
            item := v
            if item.Limit > 0 {
                page = 1
                pagesize = item.Limit
            } else {
                page = item.Page
                pagesize = item.Pagesize                
            }
            orderby = item.Order
        case Where:
            item := v

            if strings.Contains(item.Column, "_") {
                query.WriteString(" and ")
            } else {
                query.WriteString(" and a_")
            }
            query.WriteString(item.Column)
            
            if item.Compare == "in" {
                query.WriteString(" in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "not in" {
                query.WriteString(" not in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "between" {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(fmt.Sprintf(" between $%v and $%v", pos, pos + 1))
                    pos += 2
                } else {
                    query.WriteString(" between ? and ?")
                }

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(fmt.Sprintf(" $%v", pos))
                    pos++
                } else {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(" ?")
                }
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        case Custom:
             item := v

            query.WriteString(" and ")
            query.WriteString(item.Query)
        case Base:
             item := v

             baseQuery = item.Query
        }
    }

    query.WriteString(p.GroupQuery)
    
    startpage := (page - 1) * pagesize
    
    if page > 0 && pagesize > 0 {
        if orderby == "" {
            orderby = "a_id desc"
        } else {
            if !strings.Contains(orderby, "_") {                   
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "a_" + orderby
                }
            }
            
        }
        query.WriteString(" order by ")
        query.WriteString(orderby)
        if config.Database.Type == config.Postgresql {
            query.WriteString(fmt.Sprintf(" limit $%v offset $%v", pos, pos + 1))
            params = append(params, pagesize)
            params = append(params, startpage)
        } else if config.Database.Type == config.Mysql {
            query.WriteString(" limit ? offset ?")
            params = append(params, pagesize)
            params = append(params, startpage)
        } else if config.Database.Type == config.Sqlserver {
            query.WriteString("OFFSET ? ROWS FETCH NEXT ? ROWS ONLY")
            params = append(params, startpage)
            params = append(params, pagesize)
        }
    } else {
        if orderby == "" {
            orderby = "a_id"
        } else {
            if !strings.Contains(orderby, "_") {
                if strings.ToUpper(orderby) != "RAND()" {
                    orderby = "a_" + orderby
                }
            }
        }
        query.WriteString(" order by ")
        query.WriteString(orderby)
    }

    rows, err := p.Query(baseQuery + query.String(), params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
        items := make([]Answer, 0)
        return items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}


func (p *AnswerManager) FindByUid(uid int64, args ...interface{}) []Answer {
    rets := make([]interface{}, 0)
    rets = append(rets, args...)

    if uid != 0 { 
        rets = append(rets, Where{Column:"uid", Value:uid, Compare:"="})
     }
    
    return p.Find(rets)
}

func (p *AnswerManager) FindByQid(qid int64, args ...interface{}) []Answer {
    rets := make([]interface{}, 0)
    rets = append(rets, args...)

    if qid != 0 { 
        rets = append(rets, Where{Column:"qid", Value:qid, Compare:"="})
     }
    
    return p.Find(rets)
}




func (p *AnswerManager) GroupBy(name string, args []interface{}) []Groupby {
    if !p.Conn.IsConnect() {
        var items []Groupby
        return items
    }

    var params []interface{}
    baseQuery := p.GetQueryGroup(name)
    var query strings.Builder
    pos := 1

    for _, arg := range args {
        switch v := arg.(type) {
        case Where:
            item := v

            if strings.Contains(item.Column, "_") {
                query.WriteString(" and ")
            } else {
                query.WriteString(" and a_")
            }
            query.WriteString(item.Column)
            
            if item.Compare == "in" {
                query.WriteString(" in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "not in" {
                query.WriteString(" not in (")
                query.WriteString(strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]"))
                query.WriteString(")")
            } else if item.Compare == "between" {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(fmt.Sprintf(" between $%v and $%v", pos, pos + 1))
                    pos += 2
                } else {
                    query.WriteString(" between ? and ?")
                }

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                if config.Database.Type == config.Postgresql {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(fmt.Sprintf(" $%v", pos))
                    pos++
                } else {
                    query.WriteString(" ")
                    query.WriteString(item.Compare)
                    query.WriteString(" ?")
                }
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        case Custom:
             item := v

            query.WriteString(" and ")
            query.WriteString(item.Query)
        case Base:
             item := v

             baseQuery = item.Query
        }
    }
    
    query.WriteString(" group by a_")
    query.WriteString(name)

    rows, err := p.Query(baseQuery + query.String(), params...)

    if err != nil {
       if p.Log {
          log.Error().Str("error", err.Error()).Msg("SQL")
       }
        var items []Groupby
        return items
    }

    defer rows.Close()

    var items []Groupby

    for rows.Next() {
        var item Groupby
        err := rows.Scan(&item.Value, &item.Count)
        if err != nil {
           if p.Log {
                log.Error().Str("error", err.Error()).Msg("SQL")
           }
           break
        }

        items = append(items, item)
    }

    return items
}
