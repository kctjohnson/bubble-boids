# Bubble Boids!

This is just a little project I made in the hopes of learning more about Go, as well as familiarize myself
a bit more with the Charm project!

## Usage

Usage: `go run main.go`
In the main.go file, theres a const variable `servermode`. Set that to true to make this
run as a wish app.

Under ./internal/boid is a file called constants.go, and in it are a few key pieces to how the boids interact
with each other.

- BoidCount: The number of boids in the simulation
- MaxForce: The max amount of force a single part of the boid logic can exert (alignment, cohesion, separation)
- MaxSpeed: The max speed that a boid can move
- Perception: How far out the boid can see to other boids
- ScatterCounterCap: How frequently the boids scatter in different directions. It's based on fps, so 600 would be 60 seconds at 60FPS

Under ./cmd/cli is another file called constants.go, this file relates more to terminal specific constants

- TermRatio: Height to width ratio. In my terminal, the characters are each ~3 times taller than they are wide.
  This is used to make sure that the boids fill the entire terminal as well as move at the same speed up and down, and left and right
- VirtScreenWidth: How wide the virtual screen is. The boids are rendered in a virtual space, then their positions
  are squished down into the terminal space

## Credits

Best to give credit where credit is due:

- [Charm](https://charm.sh/)
- [Flapioca](https://github.com/kbrgl/flapioca) (Took the TickMsg concept from here, super cool project!)
- [The Coding Train Flocking Simulation](https://www.youtube.com/watch?v=mhjuuHl6qHM) (Where I first learned about boids, and the math behind them)

## TODO

[x] ~Separate screen and virtual screen~
[ ] Implement view ratio (Zooming in and out)
[ ] Boid directional characters (maybe use unicode?) (> v < ^)
[x] ~Optimize boid mechanisms so it's not iterating all boids 3 times per boid~
[ ] Use go concurrency to increase performance
[ ] Bubble tea UI! Help text, border box around rendered screen, etc
[x] ~Make this thing a wish app! I want to ssh into boids.internet and get a terminal with boids.~
