package gisp

import (
	"fmt"
	"reflect"

	p "github.com/Dwarfartisan/goparsec2"
)

// Axiom 是基本的 LISP 公理实现，尽可能贴近原始的 LISP 公理描述，但是部分实现对实际的 golang
// 环境做了妥协
var Axiom = Toolkit{
	Meta: map[string]interface{}{
		"name":     "axiom",
		"category": "package",
	},
	Content: map[string]interface{}{
		"quote": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Quote Args Error: expect only one arg but %v", args)
			}
			return Q(Q(args[0])), nil
		}),
		"var": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
			st := p.NewBasicState(args)
			_, err := TypeAs(ATOM).Then(p.Choice(p.Try(p.EOF), p.P(p.One).Then(p.EOF)))(&st)
			if err != nil {
				return nil, err
			}
			first := args[0].(Atom)

			slot := VarSlot(first.Type)
			if len(args) == 1 {
				err := env.Defvar(first.Name, slot)
				return Q(nil), err
			}
			val, err := Eval(env, args[1])
			slot.Set(val)
			err = env.Defvar(first.Name, slot)
			return Q(val), err
		}),
		"set": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
			st := p.NewBasicState(args)
			_, err := p.Choice(p.Try(TypeAs(ATOM)), TypeAs(QUOTE)).Then(p.One).Then(p.EOF)(&st)
			if err != nil {
				return nil, err
			}
			val, err := Eval(env, args[1])
			if err != nil {
				return nil, err
			}
			ret, err := set(env, args[0], val)
			if err != nil {
				return nil, err
			}
			return Q(ret), err
		}),
		"equal": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("args error: equal need two args but only",
					args)
			}
			return func(env Env) (interface{}, error) {
				return reflect.DeepEqual(args[0], args[1]), nil
			}, nil
		}),
		"cond": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
			cases := args[0].([]interface{})
			l := len(args)
			var els interface{}
			if l > 1 {
				els = args[1]
			} else {
				els = nil
			}

			for _, b := range cases { // FIXME: need a else
				branch := b.([]interface{})
				cond := branch[0].(List)
				result, err := Eval(env, cond)
				if err != nil {
					return nil, err
				}
				if ok := result.(bool); ok {
					return func(env Env) (interface{}, error) {
						return Eval(env, branch[1])
					}, nil
				}
			}
			// else branch
			if els != nil {
				return func(env Env) (interface{}, error) {
					return Eval(env, els)
				}, nil
			}
			return nil, nil
		}),
		"car": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
			if lisp, ok := args[0].(List); ok {
				return Q(lisp[0]).Eval, nil
			}
			return nil, fmt.Errorf("car args error: excpet a list but %v", args)
		}),
		"cdr": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
			if lisp, ok := args[0].(List); ok {
				return Q(lisp[1:]).Eval, nil
			}
			return nil, fmt.Errorf("car args error: excpet a list but %v", args)
		}),
		// atom while true both lisp atom or go value
		"atom": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
			arg := args[0]
			if l, ok := arg.(List); ok {
				return Q(len(l) == 0), nil
			}
			return Q(true), nil
		}),
		// 照搬 cons 运算符对于 golang 嵌入没有足够的收益，这里的 concat 是一个 cons 的变形，
		// 它总是返回包含所有参数的 List 。
		"concat": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
			return func(env Env) (interface{}, error) {
				return Q(List(args)).Eval, nil
			}, nil
		}),
	},
}

func set(env Env, slot, arg interface{}) (interface{}, error) {
	switch setter := slot.(type) {
	case Atom:
		err := env.Setvar(setter.Name, arg)
		if err == nil {
			return nil, err
		}
		return arg, nil
	case Bracket:
		return setter.SetItemBy(env, arg)
	case List:
		s, err := Eval(env, setter)
		if err != nil {
			return nil, err
		}
		return set(env, s, arg)
	case Quote:
		s, err := Eval(env, setter)
		if err != nil {
			return nil, err
		}
		return set(env, s, arg)
	default:
		return arg, fmt.Errorf("set error: set %v(%v) as %v is invalid",
			slot, reflect.TypeOf(slot), arg)
	}
}
