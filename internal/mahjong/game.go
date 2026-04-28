package mahjong

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
)

var fullSet = [136]rune{
	'🀀', '🀁', '🀂',
	'🀃', '🀄', '🀅', '🀆',
	'🀇', '🀈', '🀉', '🀊', '🀋', '🀌', '🀍', '🀎', '🀏',
	'🀐', '🀑', '🀒', '🀓', '🀔', '🀕', '🀖', '🀗', '🀘',
	'🀙', '🀚', '🀛', '🀜', '🀝', '🀞', '🀟', '🀠', '🀡',
	'🀀', '🀁', '🀂',
	'🀃', '🀄', '🀅', '🀆',
	'🀇', '🀈', '🀉', '🀊', '🀋', '🀌', '🀍', '🀎', '🀏',
	'🀐', '🀑', '🀒', '🀓', '🀔', '🀕', '🀖', '🀗', '🀘',
	'🀙', '🀚', '🀛', '🀜', '🀝', '🀞', '🀟', '🀠', '🀡',
	'🀀', '🀁', '🀂',
	'🀃', '🀄', '🀅', '🀆',
	'🀇', '🀈', '🀉', '🀊', '🀋', '🀌', '🀍', '🀎', '🀏',
	'🀐', '🀑', '🀒', '🀓', '🀔', '🀕', '🀖', '🀗', '🀘',
	'🀙', '🀚', '🀛', '🀜', '🀝', '🀞', '🀟', '🀠', '🀡',
	'🀀', '🀁', '🀂',
	'🀃', '🀄', '🀅', '🀆',
	'🀇', '🀈', '🀉', '🀊', '🀋', '🀌', '🀍', '🀎', '🀏',
	'🀐', '🀑', '🀒', '🀓', '🀔', '🀕', '🀖', '🀗', '🀘',
	'🀙', '🀚', '🀛', '🀜', '🀝', '🀞', '🀟', '🀠', '🀡',
}

type Game struct {
	Wall    Set
	Players []Player
}

func NewGame(playerCount int) *Game {

	players := make([]Player, playerCount)

	for i := range playerCount {
		players[i].ID = i + 1
	}

	g := Game{
		Wall:    NewWall(),
		Players: players,
	}

	g.Deal()

	return &g
}

func (g *Game) Deal() {
	const handSize = 13

	for i := range g.Players {
		g.Players[i].Hand = make([]rune, handSize)
		copy(g.Players[i].Hand, g.Wall.Stones[:handSize])
		g.Wall.Stones = g.Wall.Stones[handSize:]
	}

}

func (g *Game) Print() {
	fmt.Println(g)
}

func (g *Game) String() string {
	sb := strings.Builder{}
	sb.WriteString(g.Wall.String() + "\n")
	for _, p := range g.Players {
		sb.WriteString(p.String() + "\n")
	}
	return sb.String()
}

type Set struct {
	Stones []rune
}

func NewWall() Set {
	s := Set{Stones: make([]rune, len(fullSet))}
	copy(s.Stones, fullSet[:])
	rand.Shuffle(len(s.Stones), func(i, j int) {
		s.Stones[i], s.Stones[j] = s.Stones[j], s.Stones[i]
	})
	return s
}

func (s Set) String() string {
	return fmt.Sprintf("Wall (%d stones): %s", len(s.Stones), string(s.Stones))
}

type Player struct {
	ID   int
	Hand []rune
}

func (p Player) String() string {
	hand := string(p.Hand)
	return fmt.Sprintf("Player %d: %s", p.ID, hand)
}

type Tile rune

type Group int

const (
	GroupMan Group = iota
	GroupPin
	GroupSou
	GroupHonor
	GroupUnknown
)

func (t Tile) Group() Group {
	switch {
	case t >= '🀇' && t <= '🀏':
		return GroupMan
	case t >= '🀙' && t <= '🀡':
		return GroupPin
	case t >= '🀐' && t <= '🀘':
		return GroupSou
	case t >= '🀀' && t <= '🀆':
		return GroupHonor
	default:
		return GroupUnknown
	}
}

func (t Tile) IsNextOf(other Tile) bool {
	g := t.Group()
	if g == GroupHonor || g == GroupUnknown {
		return false
	}
	if g != other.Group() {
		return false
	}
	return rune(t) == rune(other)+1
}

type Meld struct {
	Tiles [3]rune
	Kind  MeldKind
}

type MeldKind int

const (
	Chi MeldKind = iota
	Pon
)

func (m Meld) String() string {
	return string(m.Tiles[:])
}

func (m Meld) KindString() string {
	if m.Kind == Chi {
		return "Chi"
	}
	return "Pon"
}

func FindMelds(hand []rune) []Meld {
	sorted := make([]rune, len(hand))
	copy(sorted, hand)
	slices.Sort(sorted)

	var melds []Meld
	used := make([]bool, len(sorted))

	for i := 0; i < len(sorted)-2; i++ {
		if used[i] {
			continue
		}
		if sorted[i] == sorted[i+1] && sorted[i+1] == sorted[i+2] {
			melds = append(melds, Meld{
				Tiles: [3]rune{sorted[i], sorted[i+1], sorted[i+2]},
				Kind:  Pon,
			})
			used[i], used[i+1], used[i+2] = true, true, true
			i += 2
		}
	}

	for i := range sorted {
		if used[i] || Tile(sorted[i]).Group() == GroupHonor {
			continue
		}
		j := findFirstUnused(sorted, used, sorted[i]+1)
		if j < 0 {
			continue
		}
		k := findFirstUnused(sorted, used, sorted[i]+2)
		if k < 0 {
			continue
		}
		melds = append(melds, Meld{
			Tiles: [3]rune{sorted[i], sorted[j], sorted[k]},
			Kind:  Chi,
		})
		used[i], used[j], used[k] = true, true, true
	}

	return melds
}

func FindMeldsString(hand []rune) string {
	panic("implement me")
}

func findFirstUnused(tiles []rune, used []bool, target rune) int {
	for i, t := range tiles {
		if !used[i] && t == target {
			return i
		}
	}
	return -1
}
