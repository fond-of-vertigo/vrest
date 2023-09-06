package vrest

import "testing"

func Test_makePath(t *testing.T) {
	type args struct {
		pathWithPlaceholders string
		keysAndValues        []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "1 param",
		args: args{
			pathWithPlaceholders: "/order/{Order_ID_1}",
			keysAndValues:        []string{"Order_ID_1", "1"},
		},
		want: "/order/1",
	}, {
		name: "2 params",
		args: args{
			pathWithPlaceholders: "/order/{order_id}/item/{item_id}",
			keysAndValues:        []string{"order_id", "1", "item_id", "2"},
		},
		want: "/order/1/item/2",
	}, {
		name: "No params",
		args: args{
			pathWithPlaceholders: "/order",
		},
		want: "/order",
	}, {
		name: "No params at all",
		want: "",
	}, {
		name: "Invalid number of keysAndValues",
		args: args{
			pathWithPlaceholders: "/order/{order_id}",
			keysAndValues:        []string{"order_id"},
		},
		want: "/order/{order_id}",
	}, {
		name: "Invalid key in keysAndValues",
		args: args{
			pathWithPlaceholders: "/order/{order_id}",
			keysAndValues:        []string{"xyz", "1"},
		},
		want: "/order/{order_id}",
	}, {
		name: "Invalid placeholder",
		args: args{
			pathWithPlaceholders: "/order/{!*%}",
		},
		want: "/order/{!*%}",
	}, {
		name: "Empty placeholder",
		args: args{
			pathWithPlaceholders: "/order/{}",
		},
		want: "/order/{}",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			if len(tt.args.keysAndValues) > 0 {
				got = makePath(tt.args.pathWithPlaceholders, tt.args.keysAndValues...)
			} else {
				got = makePath(tt.args.pathWithPlaceholders)
			}
			if got != tt.want {
				t.Fatalf("makePath() failed:\ngot : \"%s\"\nwant: \"%s\"", got, tt.want)
			}
		})
	}
}
