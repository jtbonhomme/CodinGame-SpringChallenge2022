package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

const (
	MonsterType  int = 0
	MyHeroType   int = 1
	OpHeroType   int = 2
	DangerRadius int = 5000
)

type Entity struct {
	id           int
	_type        int
	x, y         int
	shieldLife   int
	isControlled int
	health       int
	vx, vy       int
	nearBase     int
	threatFor    int
	risk         int
}

type Monsters []Entity
type Heroes []Entity

var baseX, baseY int

func main() {
	// baseX: The corner of the map representing your base
	fmt.Scan(&baseX, &baseY)
	// heroesPerPlayer: Always 3
	var heroesPerPlayer int
	fmt.Scan(&heroesPerPlayer)

	for {
		fmt.Fprintf(os.Stderr, "baseX baseY : %d %d\n", baseX, baseY)
		for i := 0; i < 2; i++ {
			// health: Each player's base health
			// mana: Ignore in the first league; Spend ten mana to cast a spell
			var health, mana int
			fmt.Scan(&health, &mana)
			fmt.Fprintf(os.Stderr, "hero %d - health %d mana %d\n", i, health, mana)
		}
		// entityCount: Amount of hero and monsters you can see
		var entityCount int
		fmt.Scan(&entityCount)
		fmt.Fprintf(os.Stderr, "entityCount : %d\n", entityCount)
		monsters := Monsters{}
		myHeroes := Heroes{}
		oppHeroes := Heroes{}

		for i := 0; i < entityCount; i++ {
			// id: Unique identifier
			// _type: 0=monster, 1=your hero, 2=opponent hero
			// x: Position of this entity
			// shieldLife: Ignore for this league; Count down until shield spell fades
			// isControlled: Ignore for this league; Equals 1 when this entity is under a control spell
			// health: Remaining health of this monster
			// vx: Trajectory of this monster
			// nearBase: 0=monster with no target yet, 1=monster targeting a base
			// threatFor: Given this monster's trajectory, is it a threat to 1=your base, 2=your opponent's base, 0=neither
			var id, _type, x, y, shieldLife, isControlled, health, vx, vy, nearBase, threatFor int
			fmt.Scan(&id, &_type, &x, &y, &shieldLife, &isControlled, &health, &vx, &vy, &nearBase, &threatFor)
			entity := Entity{id, _type, x, y, shieldLife, isControlled, health, vx, vy, nearBase, threatFor, 0}

			switch _type {
			case MonsterType:
				// risk computation
				entity.riskUpdate()
				monsters = append(monsters, entity)
				break
			case MyHeroType:
				myHeroes = append(myHeroes, entity)
				break
			case OpHeroType:
				oppHeroes = append(oppHeroes, entity)
				break
			}
		}
		fmt.Fprintf(os.Stderr, "monsters : %s\n", monsters.String())
		// In the first league: MOVE <x> <y> | WAIT; In later leagues: | SPELL <spellParams>;
		targets := monsters.Targets()
		fmt.Fprintf(os.Stderr, "targets: %s\n", targets.String())

		// todo: assign a target to the closest hero, not the closest target to each hero
		for i := 0; i < heroesPerPlayer; i++ {
			var minDistance = 20000
			indiceTarget := -1
			for j := 0; j < len(targets); j++ {
				if int(myHeroes[i].DistanceFrom(targets[j].x, targets[j].y)) < minDistance {
					minDistance = int(myHeroes[i].DistanceFrom(targets[j].x, targets[j].y))
					indiceTarget = j
				}
			}
			if indiceTarget != -1 {
				target := targets[indiceTarget]
				fmt.Printf("MOVE %d %d [hero %d]\n", target.x, target.y, i)
				// Remove the element at index i from a.
				targets[indiceTarget] = targets[len(targets)-1] // Copy last element to index i.
				targets = targets[:len(targets)-1]              // Truncate slice.
			} else {
				fmt.Println("WAIT")
			}
		}
	}
}

//const (
//	Zone0Radius int = 3000
//	Zone1Radius int = 5000
//	Zone2Radius int = 7000
//)

//func zoneRadius(i int) int {
//	if i >= 2 {
//		return Zone2Radius
//	}
//	if i <= 0 {
//		return Zone0Radius
//	}
//	return Zone1Radius
//}

//func (m Monsters) IsZoneEmpty(i int) bool {
//	for _, monster := range m {
//		if Entity(monster).DistanceFrom(baseX, baseY) < zoneRadius(i) {
//			return false
//		}
//	}
//	return true
//}

//func (e Entity) CurrentZone() int {
//	for i := 0; i < 3; i++ {
//		if e.DistanceFrom(baseX, baseY) < zoneRadius(i) {
//			return i
//		}
//	}
//	return 2
//}
//
//func (e Entity) DefendZone(i int, monsters Monsters) {
//	order := "WAIT"
//	for _, m := range monsters {
//		if m.DistanceFrom(baseX, baseY) < zoneRadius(i) {
//			order = fmt.Sprintf("MOVE %d %d [hero %d]", m.x, m.y, i)
//		}
//	}
//	fmt.Println(order)
//}

// By is the type of a "less" function that defines the ordering of its Entitys arguments.
type By func(e1, e2 *Entity) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(entities []Entity) {
	ps := &entitySorter{
		entities: entities,
		by:       by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// entitySorter joins a By function and a slice of Entities to be sorted.
type entitySorter struct {
	entities []Entity
	by       func(p1, p2 *Entity) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *entitySorter) Len() int {
	return len(s.entities)
}

// Swap is part of sort.Interface.
func (s *entitySorter) Swap(i, j int) {
	s.entities[i], s.entities[j] = s.entities[j], s.entities[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *entitySorter) Less(i, j int) bool {
	return s.by(&s.entities[i], &s.entities[j])
}

func (e Entity) DistanceFrom(pointX, pointY int) float64 {
	return math.Floor(math.Sqrt(float64(pointY-e.y)*float64(pointY-e.y) + float64(pointX-e.x)*float64(pointX-e.x)))
}

func (e Entity) DistanceFromBase() float64 {
	return e.DistanceFrom(baseX, baseY)
}

func (e Entity) velocity() float64 {
	return math.Floor(math.Sqrt(float64(e.vy)*float64(e.vy) + float64(e.vx)*float64(e.vx)))
}

func (e Entity) thetaMin() float64 {
	dy := float64(baseY - DangerRadius - e.y)
	dx := float64(baseX - e.x)
	t := math.Atan(dy / dx)
	if dx < 0 && dy > 0 {
		t += math.Pi
	}
	if dx < 0 && dy < 0 {
		t -= math.Pi
	}
	return t
}

func (e Entity) thetaMax() float64 {
	dy := float64(baseY - e.y)
	dx := float64(baseX - DangerRadius - e.x)
	t := math.Atan(dy / dx)
	if dx < 0 && dy > 0 {
		t += math.Pi
	}
	if dx < 0 && dy < 0 {
		t -= math.Pi
	}
	return t
}

func (e Entity) theta() float64 {
	t := math.Atan(float64(e.vy) / float64(e.vx))
	if e.vx < 0 && e.vy > 0 {
		t += math.Pi
	}
	if e.vx < 0 && e.vy < 0 {
		t -= math.Pi
	}
	return t
}

func (e *Entity) riskUpdate() {
	e.risk = 20000 - int(e.DistanceFromBase())
	if e.theta() > e.thetaMin() && e.theta() < e.thetaMax() {
		e.risk += 20000
	}
	if e.nearBase == 1 && e.DistanceFromBase() <= 5000.0 {
		e.risk += 40000
	}
}

func (m Monsters) String() string {
	res := "\n"
	for _, e := range m {
		res += "\n" + e.String()
	}
	return res
}

func (e Entity) String() string {
	return fmt.Sprintf("[%d] health %d - pos (%d,%d) - velocity (%d, %d) - theta %.2f (%.2f, %.2f) - nearBase %d - risk %d", e.id, e.health, e.x, e.y, e.vx, e.vy, e.theta(), e.thetaMin(), e.thetaMax(), e.nearBase, e.risk)
}

func (m Monsters) Targets() Monsters {
	distance := func(e1, e2 *Entity) bool {
		return e1.risk < e2.risk
	}
	decreasingDistance := func(e1, e2 *Entity) bool {
		return distance(e2, e1)
	}

	monsters := m
	By(decreasingDistance).Sort(monsters)
	fmt.Fprintf(os.Stderr, "sorted monsters: %s\n", monsters.String())

	targets := []Entity{}
	for i := 0; i < 3 && i < len(monsters); i++ {
		targets = append(targets, monsters[i])
	}
	return targets
}
