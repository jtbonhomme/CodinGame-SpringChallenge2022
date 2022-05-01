# CodinGame-SpringChallenge2022
My AI for [Codingame Spring Challenge 2022](https://www.codingame.com/contests/spring-challenge-2022) in Go

## Strategy

### Wood ligue 2

For this round, strategy consists in ranking monsters regarding to their distance from the base (risk is higher if distance is smaller).
To make sure we focus on the most dangerous one, the heuristic adds an extra risk if the monsters move to the base.
An other extra risk is added if the monsters are located closer than 5000 units from the base.


