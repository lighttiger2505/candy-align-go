package main

import (
	"reflect"
	"testing"
)

func Test_splitToTable(t *testing.T) {
	type args struct {
		str       string
		delimiter string
	}
	tests := []struct {
		name  string
		args  args
		want  [][]string
		want1 int
	}{
		{
			"single whitespace separate",
			args{
				str:       "foo bar foobar\nfoobar bar foo",
				delimiter: "",
			},
			[][]string{
				{"foo", "bar", "foobar"},
				{"foobar", "bar", "foo"},
			},
			3,
		},
		{
			"multiple whitespace separate",
			args{
				str:       "foo  bar   foobar\nfoobar   bar  foo",
				delimiter: "",
			},
			[][]string{
				{"foo", "bar", "foobar"},
				{"foobar", "bar", "foo"},
			},
			3,
		},
		{
			"tab separate",
			args{
				str:       "foo\tbar\tfoobar\nfoobar\tbar\tfoo",
				delimiter: "",
			},
			[][]string{
				{"foo", "bar", "foobar"},
				{"foobar", "bar", "foo"},
			},
			3,
		},
		{
			"specific delimiter",
			args{
				str:       "foo,bar,foobar\nfoobar,bar,foo",
				delimiter: ",",
			},
			[][]string{
				{"foo", "bar", "foobar"},
				{"foobar", "bar", "foo"},
			},
			3,
		},
		{
			"specific delimiter with whitespace",
			args{
				str:       " foo,bar , foobar \n  foobar,bar  ,  foo  ",
				delimiter: ",",
			},
			[][]string{
				{"foo", "bar", "foobar"},
				{"foobar", "bar", "foo"},
			},
			3,
		},
		{
			"Contains a blank line",
			args{
				str:       "foo bar foobar\nfoobar bar foo\n\n\n",
				delimiter: "",
			},
			[][]string{
				{"foo", "bar", "foobar"},
				{"foobar", "bar", "foo"},
			},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitToTable(tt.args.str, tt.args.delimiter)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitToTable() \ngot: %v\nwant %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitToTable() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_countFields(t *testing.T) {
	type args struct {
		table      [][]string
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
			"upper count multi byte charcter",
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
			"lower count multi byte charcter",
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
			if got := countFields(tt.args.table, tt.args.columnSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_padFields(t *testing.T) {
	type args struct {
		table  [][]string
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
			if got := padFields(tt.args.table, tt.args.counts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("padFields() = \ngot :%v\nwant:%v", got, tt.want)
			}
		})
	}
}

func Benchmark_padFields(b *testing.B) {
	table := [][]string{
		{"f", "ba", "foo"},
		{"あ", "あい", "あいう"},
		{"foo", "barr", "fooba"},
		{"あい", "あいう", "あいうえ"},
	}
	counts := []int{4, 6, 8}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		padFields(table, counts)
	}
}

func Test_trancateProtrudeString(t *testing.T) {
	type args struct {
		table  [][]string
		limits []int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			"ok",
			args{
				[][]string{
					{"f   ", "ba    ", "foo     "},
					{"あ  ", "あい  ", "あいう  "},
					{"foo ", "barr  ", "fooba   "},
					{"あい", "あいう", "あいうえ"},
				},
				[]int{2, 4, 6},
			},
			[][]string{
				{"f ", "ba  ", "foo   "},
				{"あ", "あい", "あいう"},
				{"fo", "barr", "fooba "},
				{"あ", "あい", "あいう"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trancateProtrudeString(tt.args.table, tt.args.limits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trancateProtrudeString() = \ngot: %v\nwant %v", got, tt.want)
			}
		})
	}
}

func Test_parceLimits(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "no white space",
			args: args{
				str: "1,2,3,4",
			},
			want:    []int{1, 2, 3, 4},
			wantErr: false,
		},
		{
			name: "white space separator after",
			args: args{
				str: "1 ,2 ,3 ,4",
			},
			want:    []int{1, 2, 3, 4},
			wantErr: false,
		},
		{
			name: "white space separator before",
			args: args{
				str: "1, 2, 3, 4",
			},
			want:    []int{1, 2, 3, 4},
			wantErr: false,
		},
		{
			name: "multiple white space",
			args: args{
				str: "1, 2,  3,   4",
			},
			want:    []int{1, 2, 3, 4},
			wantErr: false,
		},
		{
			name: "contains charcter",
			args: args{
				str: "1,a,3,4",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parceWidthFlag(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("parceWidthFlag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parceWidthFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createDrawLines(t *testing.T) {
	type args struct {
		table     [][]string
		delimiter string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "default delimiter",
			args: args{
				table: [][]string{
					{"foo", "bar", "foobar"},
					{"foobar", "bar", "foo"},
				},
				delimiter: "",
			},
			want: []string{
				"foo\tbar\tfoobar",
				"foobar\tbar\tfoo",
			},
		},
		{
			name: "specific delimiter",
			args: args{
				table: [][]string{
					{"foo", "bar", "foobar"},
					{"foobar", "bar", "foo"},
				},
				delimiter: ",",
			},
			want: []string{
				"foo,bar,foobar",
				"foobar,bar,foo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createDrawLines(tt.args.table, tt.args.delimiter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDrawLines() = %v, want %v", got, tt.want)
			}
		})
	}
}
