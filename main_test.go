package main

import (
	"reflect"
	"testing"
)

func Test_toSheetString(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name  string
		args  args
		want  [][]string
		want1 int
	}{
		{
			"single whitespace",
			args{"foo bar foobar\nfoobar bar foo"},
			[][]string{
				{"foo", "bar", "foobar"},
				{"foobar", "bar", "foo"},
			},
			3,
		},
		{
			"multiple whitespace",
			args{"foo  bar   foobar\nfoobar   bar  foo"},
			[][]string{
				{"foo", "bar", "foobar"},
				{"foobar", "bar", "foo"},
			},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := toSheetString(tt.args.val)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toSheetString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("toSheetString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_countColumn(t *testing.T) {
	type args struct {
		sheet      [][]string
		columnSize int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"upper count ascii",
			args{
				[][]string{
					{"f", "ba", "foo"},
					{"fo", "bar", "foob"},
					{"foo", "barr", "fooba"},
					{"fooo", "barrr", "foobar"},
				},
				3,
			},
			[]int{4, 5, 6},
		},
		{
			"upper count multi byte word",
			args{
				[][]string{
					{"あ", "あい", "あいう"},
					{"あい", "あいう", "あいうえ"},
					{"あいう", "あいうえ", "あいうえお"},
				},
				3,
			},
			[]int{6, 8, 10},
		},
		{
			"lower count ascii",
			args{
				[][]string{
					{"fooo", "barrr", "foobar"},
					{"foo", "barr", "fooba"},
					{"fo", "bar", "foob"},
					{"f", "ba", "foo"},
				},
				3,
			},
			[]int{4, 5, 6},
		},
		{
			"lower count multi byte word",
			args{
				[][]string{
					{"あいう", "あいうえ", "あいうえお"},
					{"あい", "あいう", "あいうえ"},
					{"あ", "あい", "あいう"},
				},
				3,
			},
			[]int{6, 8, 10},
		},
		{
			"upper count mix",
			args{
				[][]string{
					{"f", "ba", "foo"},
					{"あ", "あい", "あいう"},
					{"foo", "barr", "fooba"},
					{"あい", "あいう", "あいうえ"},
				},
				3,
			},
			[]int{4, 6, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countColumn(tt.args.sheet, tt.args.columnSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paddingSheet(t *testing.T) {
	type args struct {
		sheet  [][]string
		counts []int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			"ascii only",
			args{
				[][]string{
					{"fooo", "barrr", "foobar"},
					{"foo", "barr", "fooba"},
					{"fo", "bar", "foob"},
					{"f", "ba", "foo"},
				},
				[]int{4, 5, 6},
			},
			[][]string{
				{"fooo", "barrr", "foobar"},
				{"foo ", "barr ", "fooba "},
				{"fo  ", "bar  ", "foob  "},
				{"f   ", "ba   ", "foo   "},
			},
		},
		{
			"multi byte only",
			args{
				[][]string{
					{"あ", "あい", "あいう"},
					{"あい", "あいう", "あいうえ"},
					{"あいう", "あいうえ", "あいうえお"},
				},
				[]int{6, 8, 10},
			},
			[][]string{
				{"あ    ", "あい    ", "あいう    "},
				{"あい  ", "あいう  ", "あいうえ  "},
				{"あいう", "あいうえ", "あいうえお"},
			},
		},
		{
			"mix",
			args{
				[][]string{
					{"f", "ba", "foo"},
					{"あ", "あい", "あいう"},
					{"foo", "barr", "fooba"},
					{"あい", "あいう", "あいうえ"},
				},
				[]int{4, 6, 8},
			},
			[][]string{
				{"f   ", "ba    ", "foo     "},
				{"あ  ", "あい  ", "あいう  "},
				{"foo ", "barr  ", "fooba   "},
				{"あい", "あいう", "あいうえ"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paddingSheet(tt.args.sheet, tt.args.counts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paddingSheet() = %v, want %v", got, tt.want)
			}
		})
	}
}
