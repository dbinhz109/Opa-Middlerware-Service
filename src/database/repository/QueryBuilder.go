package repository

import (
	"go-app/src/model"
	"go-app/src/service"
	"reflect"
	"strconv"
	"strings"
)

// Logical operators for conditions
type OperatorEnum int

const (
	OP_NONE = iota
	OP_NOT
	OP_AND
	OP_OR
	OP_XOR
)

func (o OperatorEnum) String() string {
	switch o {
	case OP_NOT:
		return "NOT"
	case OP_AND:
		return "AND"
	case OP_OR:
		return "OR"
	default:
		return ""
	}
}

// Relation operators for conditions
type RelationEnum int

const (
	REL_NONE = iota
	REL_EQUAL
	REL_NOTEQUAL
	REL_GT
	REL_LT
	REL_GTE
	REL_LTE
)

// JoinType defines the type of join
type JoinTypeEnum int

const (
	JOIN_NONE = iota
	JOIN_INNER
	JOIN_OUTER
	JOIN_LEFT
	JOIN_RIGHT
	JOIN_CROSS
	JOIN_SELF
)

func (e JoinTypeEnum) String() string {
	switch e {
	case JOIN_INNER:
		return "INNER JOIN"
	case JOIN_LEFT:
		return "LEFT JOIN"
	case JOIN_RIGHT:
		return "RIGHT JOIN"
	case JOIN_OUTER:
		return "OUTER JOIN"
	case JOIN_CROSS:
		return "CROSS JOIN"
	case JOIN_SELF:
		return "SELF JOIN"
	default:
		return "JOIN"
	}
}

func (o RelationEnum) String() string {
	switch o {
	case REL_EQUAL:
		return "="
	case REL_GT:
		return ">"
	case REL_LT:
		return "<"
	case REL_GTE:
		return ">="
	case REL_LTE:
		return "<="
	default:
		return ""
	}
}

type ConditionRule struct {
	Operator OperatorEnum    `json:"op,omitempty"`
	Rules    []ConditionRule `json:"rules,omitempty"`

	Relation RelationEnum `json:"rel,omitempty"`
	FieldsValue
}

type JoinRule struct {
	Type  JoinTypeEnum           `json:"type,omitempty"`
	Left  model.EntityModel[any] `json:"left,omitempty"`
	Right model.EntityModel[any] `json:"right,omitempty"`
}

// NewRule returns a new ConditionRule
func NewRule(selector model.EntityModel[any], value model.EntityModel[any], relation RelationEnum) ConditionRule {
	return ConditionRule{
		Rules: make([]ConditionRule, 0),
		FieldsValue: FieldsValue{
			Selector: selector,
			Value:    value,
		},
		Relation: relation,
	}
}

// FieldsValue defines (sub) collection of fields by non-zero fields of Selector and corresponding values of Value
type FieldsValue struct {
	Selector model.EntityModel[any]
	Value    model.EntityModel[any]
}
type FieldValue struct {
	name  string
	value any
}

// GetFields returns a slice of FieldValue
func (v *FieldsValue) GetFields() []FieldValue {
	ret := make([]FieldValue, 0)
	if v.Selector == nil || v.Value == nil {
		return ret
	}
	v.Selector.ForEach(v.Value, ret, func(field string, value interface{}, shared any) {
		ret = append(ret, FieldValue{name: field, value: value})
	})
	return ret
}

// ForEach iterates over non-zero fields of Selector and corresponding values of Value
func (v *FieldsValue) ForEach(sharedData any, f func(fieldName string, fieldValue any, shared any)) {
	if v.Value == nil {
		return
	}
	if v.Selector == nil {
		v.Selector = v.Value
	}

	indices := make([]int, 0, MAX_FIELDS)
	idx := 0

	v.Selector.ForEach(v.Selector, nil, func(name string, value any, shared any) {
		if !reflect.ValueOf(value).IsZero() {
			indices = append(indices, idx)
		}
		idx += 1
	})
	// Collect field values as parameter array
	idx = 0
	v.Selector.ForEach(v.Value, nil, func(name string, value any, shared any) {
		if len(indices) == 0 {
			return
		}
		if indices[0] == idx {
			f(name, value, sharedData)
			indices = indices[1:]
		}
		idx += 1
	})
}

const (
	MAX_FIELDS int = 128
)

var (
	_tableAlias = []string{"t0", "t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9", "t10"}
)

type QueryBuilder struct {
	fromEntity model.EntityModel[any]
	count      bool
	whereRule  ConditionRule
	sort       []string
	offset     int64
	size       *int64
	joined     []JoinRule
	aliasMap   map[string]string
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		sort:     make([]string, 0),
		joined:   make([]JoinRule, 0),
		aliasMap: map[string]string{},
	}
}

func (q *QueryBuilder) From(entity model.EntityModel[any]) *QueryBuilder {
	q.fromEntity = entity
	return q
}

func (q *QueryBuilder) Select(entity model.EntityModel[any]) *QueryBuilder {
	q.fromEntity = entity
	return q
}

func (q *QueryBuilder) Where(rule ConditionRule) *QueryBuilder {
	q.whereRule = rule
	return q
}

func (q *QueryBuilder) OrderBy(fields []string) *QueryBuilder {
	q.sort = fields
	return q
}

func (q *QueryBuilder) Offset(offset int64) *QueryBuilder {
	q.offset = offset
	return q
}
func (q *QueryBuilder) Limit(limit int64) *QueryBuilder {
	q.size = &limit
	return q
}

func (q *QueryBuilder) Join(leftSelector model.EntityModel[any], joinType JoinTypeEnum, rightSelector model.EntityModel[any]) *QueryBuilder {
	rule := JoinRule{
		Type:  joinType,
		Left:  leftSelector,
		Right: rightSelector,
	}
	q.joined = append(q.joined, rule)
	return q
}

func (q *QueryBuilder) buildFromClause() string {
	var builder strings.Builder
	q.aliasMap = make(map[string]string)
	builder.WriteString(q.fromEntity.TName())
	builder.WriteString(" t0")
	q.aliasMap[q.fromEntity.TName()] = "t0"
	for idx, join := range q.joined {
		builder.WriteString(" ")
		builder.WriteString(join.Type.String())
		builder.WriteString(" ")
		builder.WriteString(join.Right.TName())
		builder.WriteString(" ")
		builder.WriteString(_tableAlias[idx+1])
		q.aliasMap[join.Right.TName()] = _tableAlias[idx+1]
		builder.WriteString(" ON ")
		builder.WriteString(q.aliasMap[join.Left.TName()])
		builder.WriteString(".")
		leftField := FieldsValue{Selector: join.Left, Value: join.Left}
		builder.WriteString(leftField.GetFields()[0].name)
		builder.WriteString(" = ")
		builder.WriteString(q.aliasMap[join.Right.TName()])
		builder.WriteString(".")
		rightField := FieldsValue{Selector: join.Right, Value: join.Right}
		builder.WriteString(rightField.GetFields()[0].name)
	}
	// logger.Info("buildFromClause", zap.Any("map", q.aliasMap))
	return builder.String()
}

func (q *QueryBuilder) buildSelectClause() string {
	var builder strings.Builder
	if q.count {
		builder.WriteString("count(*) total_count")
		return builder.String()
	}
	if q.fromEntity != nil {
		builder.WriteString("t0.*")
		for idx, _ := range q.joined {
			builder.WriteString(", ")
			builder.WriteString(_tableAlias[idx+1])
			builder.WriteString(".*")
		}
		return builder.String()
	}
	return ""
}

func (q *QueryBuilder) whereBinary(e model.EntityModel[any], rel RelationEnum) string {
	name := e.TName()
	alias, ok := q.aliasMap[name]
	if ok {
		name = alias
	}
	ret := QB.FormatNonNilFields(e, name+".%s"+rel.String()+"?", " AND ")
	return ret
}

func (q *QueryBuilder) whereUnary(e model.EntityModel[any], rel RelationEnum) string {
	name := e.TName()
	alias, ok := q.aliasMap[name]
	if ok {
		name = alias
	}
	ret := QB.FormatNonNilFields(e, name+".%s"+rel.String(), " AND ")
	return ret
}

func (q *QueryBuilder) buildWhereClause(r *ConditionRule) (string, []any) {
	ret := ""
	parameters := make([]any, 0)
	if r.Operator != OP_NONE {
		children := service.SliceMap(r.Rules, func(child ConditionRule) string {
			str, param := q.buildWhereClause(&child)
			parameters = append(parameters, param...)
			return str
		})
		ret = strings.Join(children, r.Operator.String())
	} else {
		if r.Relation == REL_NONE {
			return "", parameters
		}
		r.Selector.ForEach(r.Value, r, func(name string, val any, shared any) {
			parameters = append(parameters, val)
		})
		ret = q.whereBinary(r.Selector, r.Relation)
	}
	if len(ret) > 0 {
		return "(" + ret + ")", parameters
	} else {
		return "", parameters
	}
}

func (q *QueryBuilder) BuildSelect() (string, []any) {
	var builder strings.Builder
	builder.WriteString("SELECT ")
	builder.WriteString(q.buildSelectClause())
	builder.WriteString(" FROM ")
	builder.WriteString(q.buildFromClause())
	builder.WriteString(" ")
	where, parameters := q.buildWhereClause(&q.whereRule)
	if len(where) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(where)
	}
	if q.sort != nil && len(q.sort) > 0 {
		l := len(q.sort)
		builder.WriteString(" ORDER BY ")
		for idx, sort := range q.sort {
			if strings.HasSuffix(sort, "-") {
				builder.WriteString(strings.ReplaceAll(sort, "-", " DESC"))
			} else {
				builder.WriteString(sort)
				builder.WriteString(" ASC")
			}
			if l > 1 && idx < l-1 {
				builder.WriteString(", ")
			}
		}
	}
	builder.WriteString(" OFFSET ")
	builder.WriteString(strconv.FormatInt(q.offset, 10))
	if q.size != nil {
		builder.WriteString(" LIMIT ")
		builder.WriteString(strconv.FormatInt(*q.size, 10))
	}

	return builder.String(), parameters
}

func (q *QueryBuilder) BuildCreate(entity model.EntityModel[any]) (string, []any) {
	var builder strings.Builder
	parameters := make([]any, 0)
	builder.WriteString("INSERT INTO ")
	builder.WriteString(entity.TName())
	builder.WriteString(" (")
	names := make([]string, 0)
	entity.ForEach(entity, nil, func(name string, val any, shared any) {
		names = append(names, name)
	})
	if len(names) > 0 {
		builder.WriteString(strings.Join(names, ","))
	}
	builder.WriteString(") VALUES (")
	values := make([]string, 0)
	entity.ForEach(entity, nil, func(name string, val any, shared any) {
		values = append(values, "?")
		parameters = append(parameters, val)
	})
	if len(values) > 0 {
		builder.WriteString(strings.Join(values, ","))
	}
	builder.WriteString(")")
	return builder.String(), parameters
}

func (q *QueryBuilder) BuildUpdate(selector model.EntityModel[any], value model.EntityModel[any]) (string, []any) {
	var builder strings.Builder
	parameters := make([]any, 0)
	builder.WriteString("UPDATE ")
	builder.WriteString(selector.TName())
	builder.WriteString(" SET ")
	values := make([]string, 0)
	selector.ForEach(value, nil, func(name string, val any, shared any) {
		values = append(values, name+"=?")
		parameters = append(parameters, val)
	})
	if len(values) > 0 {
		builder.WriteString(strings.Join(values, ","))
	}
	// where, params := q.buildWhereClause(&q.whereRule)
	where, params := q.buildWhere(&q.whereRule)
	if len(where) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(where)
		parameters = append(parameters, params...)
	}
	return builder.String(), parameters
}

func (q *QueryBuilder) BuildDelete() (string, []any) {
	var builder strings.Builder
	builder.WriteString("DELETE FROM ")
	builder.WriteString(q.fromEntity.TName())
	where, parameters := q.buildWhereClause(&q.whereRule)
	if len(where) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(where)
	}
	return builder.String(), parameters
}

func (q *QueryBuilder) buildWhere(rule *ConditionRule) (string, []any) {
	b := strings.Builder{}
	parameters := make([]any, 0, MAX_FIELDS)

	tName := rule.Value.TName()
	alias, ok := q.aliasMap[tName]
	if ok {
		tName = alias
	}

	if rule.Operator != OP_NONE {
		// recursive call for each child
		first := true
		for _, child := range rule.Rules {
			if first {
				first = false
				b.WriteString("(")
			} else {
				b.WriteString(" ")
				b.WriteString(rule.Operator.String())
				b.WriteString(" ")
			}
			str, param := q.buildWhere(&child)
			parameters = append(parameters, param...)
			b.WriteString(str)
		}
	} else {
		first := true
		rule.ForEach(nil, func(name string, value any, shared any) {
			if first {
				first = false
				b.WriteString("(")
			} else {
				b.WriteString(" AND ")
			}
			b.WriteString(tName)
			b.WriteString(".")
			b.WriteString(name)
			if reflect.ValueOf(value).IsZero() {
				b.WriteString(" IS NULL")
			} else {
				b.WriteString("=?")
				parameters = append(parameters, value)
			}
		})
	}
	if b.Len() > 0 {
		b.WriteString(")")
	}
	return b.String(), parameters
}

// staticBuilder builds query parts from an entity model
type staticBuilder struct{}

func (qb *staticBuilder) Build() string {
	return ""
}

var QB staticBuilder

// WhereEqual build WHERE clause with equal relation from non nil entity fields |
// [WHERE id=:id AND name=:name AND string_value=:string_value]
func (qb *staticBuilder) WhereEqual(e model.EntityModel[any]) string {
	ret := qb.FormatNonNilFields(e, e.TName()+".%s=:%s", " AND ")
	return ret
}

// WhereOp build WHERE clause with given operator like =,<,>,>=,<= |
// [WHERE timestamp>=:timestamp AND string_value>=:string_value]
func (qb *staticBuilder) WhereOp(e model.EntityModel[any], op string) string {
	ret := qb.FormatNonNilFields(e, e.TName()+".%s"+op+":%s", " AND ")
	return ret
}

// WhereEqualOrder build WHERE clause to use with ordered parameters instead of named ones
// [WHERE id=? AND name=? AND string_value=?]
func (qb *staticBuilder) WhereEqualOrder(e model.EntityModel[any]) string {
	ret := qb.FormatNonNilFields(e, e.TName()+".%s=?", " AND ")
	return ret
}

func (qb *staticBuilder) WhereOpOrder(e model.EntityModel[any], op string) string {
	ret := qb.FormatNonNilFields(e, e.TName()+".%s"+op+"?", " AND ")
	return ret
}

// WhereEqual build WHERE clause checking NULL value from non nil entity fields |
// [WHERE name IS NULL AND string_value IS NULL]
func (qb *staticBuilder) WhereNull(e model.EntityModel[any]) string {
	ret := qb.FormatNonNilFields(e, e.TName()+".%s IS NULL", " AND ")
	return ret
}

func (qb *staticBuilder) WhereNotNull(e model.EntityModel[any]) string {
	ret := qb.FormatNonNilFields(e, e.TName()+".%s IS NOT NULL", " AND ")
	return ret
}

// ColumnList builds column name list from non nil entity fields |
// [id,name,string_value,number_value]
func (qb *staticBuilder) ColumnList(e model.EntityModel[any]) string {
	return qb.FormatNonNilFields(e, "%s", ",")
}

// ParamList builds named param list from non nil entity fields |
// [:id,:name,:string_value,:number_value]
func (qb *staticBuilder) NamedParamList(e model.EntityModel[any]) string {
	return qb.FormatNonNilFields(e, ":%s", ",")

}

// NamedSetValueList
// [name=?,string_value=?]
func (qb *staticBuilder) NamedSetValueList(e model.EntityModel[any]) string {
	return qb.FormatNonNilFields(e, "%s=?", ",")
}

// ValueList returns slice of non nil field values of entity
func (qb *staticBuilder) ValueList(e model.EntityModel[any]) []any {
	return qb.MapNonNilFields(e, func(name string, val any) any {
		return val
	})
}

// NamedSetNullList build set NULL value clause from non nil entity fields |
// [name=NULL,string_value=NULL,number_value=NULL]
func (qb *staticBuilder) NamedSetNullList(e model.EntityModel[any]) string {
	return qb.FormatNonNilFields(e, "%s=NULL", ",")
}

// MapNonNilFields maps non nil fields to an array []any
func (qb *staticBuilder) MapNonNilFields(e model.EntityModel[any], f func(fieldName string, fieldValue any) any) []any {
	shared := make([]any, 0, 256)
	e.ForEach(e, &shared, func(name string, value any, sharedValue any) {
		pShared, _ := sharedValue.(*[]any)
		if value != nil {
			*pShared = append(*pShared, f(name, value))
		}
	})
	return shared
}

// FormatNonNilFields builds formatted string array from non nil fields then joined by given str
func (qb *staticBuilder) FormatNonNilFields(e model.EntityModel[any], strFormat, strJoin string) string {
	fieldValueList := qb.MapNonNilFields(e, func(name string, value any) any {
		return strings.ReplaceAll(strFormat, "%s", name)
	})
	fieldStrList := make([]string, 0, len(fieldValueList))
	for _, fieldAny := range fieldValueList {
		fieldStrList = append(fieldStrList, fieldAny.(string))
	}
	if len(fieldStrList) > 0 {
		return strings.Join(fieldStrList, strJoin)
	}
	return ""
}
