package goskill

import (
	"reflect"
	"testing"
)

func TestBTFDefaults(t *testing.T) {
	got := New()
	expected := BradleyTerryFull{
		Mu:    25,
		Sigma: (25. / 3),
		Beta:  (25. / 3) / 2,
		Kappa: 0.0001,
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("got: %v, want: %v", got, expected)
	}
}

func TestBTFDefaultSkill(t *testing.T) {
	got := New().Skill()
	expected := Skill{
		Mu:    25,
		SigSq: 69.44444444444446,
	}
	if got != expected {
		t.Fatalf("got: %v, want: %v", got, expected)
	}
}

func TestCreateTeams(t *testing.T) {
	skills := []*Skill{
		&Skill{Mu: 25, SigSq: 1},
		&Skill{Mu: 13, SigSq: 12309},
	}
	expected := team{
		Players: skills,
		S: Skill{
			Mu:    38,
			SigSq: 12310,
		},
	}
	got := createTeam(skills)
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("got: %v, want: %v", got, expected)
	}
}

var BTFTests = map[string]struct {
	in   [][]Skill
	want [][]Skill
}{
	"one player game": {
		[][]Skill{[]Skill{Skill{Mu: 25, SigSq: 69.44444444444446}}},
		[][]Skill{[]Skill{Skill{Mu: 25, SigSq: 69.44444444444446}}},
	},
	"one vs one": {
		[][]Skill{
			[]Skill{Skill{Mu: 25, SigSq: 69.44444444444446}},
			[]Skill{Skill{Mu: 25, SigSq: 69.44444444444446}},
		},
		[][]Skill{
			[]Skill{Skill{Mu: 27.63523138347365, SigSq: 65.05239213865504}},
			[]Skill{Skill{Mu: 22.36476861652635, SigSq: 65.05239213865504}},
		},
	},
	"2v1": {
		[][]Skill{
			[]Skill{
				Skill{Mu: 15, SigSq: 4.1233},
				Skill{Mu: 10.2, SigSq: 9.220903}},
			[]Skill{Skill{Mu: 50, SigSq: 2.1232}},
		},
		[][]Skill{
			[]Skill{
				Skill{Mu: 15.56496998825091, SigSq: 4.118333035231164},
				Skill{Mu: 11.463437891876112, SigSq: 9.19606319864965}},
			[]Skill{Skill{Mu: 49.709081493208274, SigSq: 2.122674670057498}},
		},
	},
}

func TestBradleyTerryFull(t *testing.T) {
	bt := New()
	for name, test := range BTFTests {
		t.Run(name, func(t *testing.T) {
			teams := [][]*Skill{}
			for _, t := range test.in {
				team := []*Skill{}
				for i, _ := range t {
					team = append(team, &t[i])
				}
				teams = append(teams, team)
			}
			bt.Rank(teams)
			if !reflect.DeepEqual(test.in, test.want) {
				t.Fatalf("got: %v, want: %v", test.in, test.want)
			}
		})
	}
}

var BTFOrderedTests = map[string]struct {
	skills [][]Skill
	ranks  []int
	want   [][]Skill
}{
	"two player reversed": {
		[][]Skill{
			[]Skill{Skill{Mu: 25, SigSq: 69.44444444444446}},
			[]Skill{Skill{Mu: 25, SigSq: 69.44444444444446}},
		},
		[]int{2, 1},
		[][]Skill{
			[]Skill{Skill{Mu: 22.36476861652635, SigSq: 65.05239213865504}},
			[]Skill{Skill{Mu: 27.63523138347365, SigSq: 65.05239213865504}},
		},
	},
	"2v1 reversed": {
		[][]Skill{
			[]Skill{Skill{Mu: 50, SigSq: 2.1232}},
			[]Skill{
				Skill{Mu: 15, SigSq: 4.1233},
				Skill{Mu: 10.2, SigSq: 9.220903}},
		},
		[]int{10, 1},
		[][]Skill{
			[]Skill{Skill{Mu: 49.709081493208274, SigSq: 2.122674670057498}},
			[]Skill{
				Skill{Mu: 15.56496998825091, SigSq: 4.118333035231164},
				Skill{Mu: 11.463437891876112, SigSq: 9.19606319864965}},
		},
	},
	"two player tie": {
		[][]Skill{
			[]Skill{Skill{Mu: 25, SigSq: 69.44444444444446}},
			[]Skill{Skill{Mu: 25, SigSq: 69.44444444444446}},
		},
		[]int{1, 1},
		[][]Skill{
			[]Skill{Skill{Mu: 25, SigSq: 65.05239213865504}},
			[]Skill{Skill{Mu: 25, SigSq: 65.05239213865504}},
		},
	},
}

func TestBradleyTerryFullOrdered(t *testing.T) {
	bt := New()
	for name, test := range BTFOrderedTests {
		t.Run(name, func(t *testing.T) {
			teams := [][]*Skill{}
			for _, t := range test.skills {
				team := []*Skill{}
				for i, _ := range t {
					team = append(team, &t[i])
				}
				teams = append(teams, team)
			}
			bt.RankOrdered(teams, test.ranks)
			if !reflect.DeepEqual(test.skills, test.want) {
				t.Fatalf("got: %v, want: %v", test.skills, test.want)
			}
		})
	}
}
