package goskill

import (
	"math"
)

type BradleyTerryFull struct {
	Mu    float64
	Sigma float64
	Beta  float64
	Kappa float64
}

type Skill struct {
	Mu    float64
	SigSq float64
}

type team struct {
	Players []*Skill
	S       Skill
}

func (bt BradleyTerryFull) ciq(i, q Skill) float64 {
	return math.Sqrt(i.SigSq + q.SigSq + 2*bt.Beta*bt.Beta)
}

func (bt BradleyTerryFull) piq(i, q Skill, ciq float64) float64 {
	return 1 / (1 + math.Exp((q.Mu-i.Mu)/ciq))
}

func (bt BradleyTerryFull) rank(teams []team, ranks []int) {
	for i, ti := range teams {
		var omega float64
		var delta float64
		for q, tq := range teams {
			if q == i {
				continue
			}
			ciq := bt.ciq(ti.S, tq.S)
			piq := bt.piq(ti.S, tq.S, ciq)
			s := 0.
			if ranks[q] > ranks[i] {
				s = 1.
			} else if ranks[q] == ranks[i] {
				s = 0.5
			}
			omega += ti.S.SigSq / ciq * (s - piq)
			// TODO: make gamma configurable
			gamma := math.Sqrt(ti.S.SigSq) / ciq
			delta += gamma * ti.S.SigSq / ciq / ciq * piq * (1 - piq)
		}
		for _, p := range ti.Players {
			p.Mu += p.SigSq / ti.S.SigSq * omega
			p.SigSq *= math.Max((1 - p.SigSq/ti.S.SigSq*delta), bt.Kappa)
		}
	}
	return
}

func createTeam(skills []*Skill) team {
	t := team{Players: skills}
	for _, s := range skills {
		t.S.Mu += s.Mu
		t.S.SigSq += s.SigSq
	}
	return t
}

func New() BradleyTerryFull {
	return BradleyTerryFull{
		Mu:    25,
		Sigma: 25. / 3,
		Beta:  25. / 3 / 2,
		Kappa: 0.0001,
	}
}

func (bt BradleyTerryFull) Rank(skills [][]*Skill) {
	teams := []team{}
	ranks := []int{}
	for i, t := range skills {
		teams = append(teams, createTeam(t))
		ranks = append(ranks, i)
	}
	bt.rank(teams, ranks)
	return
}

func (bt BradleyTerryFull) RankOrdered(skills [][]*Skill, ranks []int) {
	teams := []team{}
	for _, t := range skills {
		teams = append(teams, createTeam(t))
	}
	bt.rank(teams, ranks)
	return
}

func (bt BradleyTerryFull) Skill() Skill {
	return Skill{
		Mu:    bt.Mu,
		SigSq: bt.Sigma * bt.Sigma,
	}
}

func (bt BradleyTerryFull) WinProbability(a, b []*Skill) float64 {
	teamA := createTeam(a)
	teamB := createTeam(b)
	return bt.piq(teamA.S, teamB.S, bt.ciq(teamA.S, teamB.S))
}
