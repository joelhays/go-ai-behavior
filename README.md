go-ai-behavior
================

Layered AI Behavior demo written in Go.

Supported actor behaviors:
* Avoid - maintain a minimum distance from the target
* Seek - follow the target
* Wander - move toward a new random direction periodically
* Constant - move in a set direction
* Rotate - change orientation
* Keyboard - enable WASD keyboard control

The scenario is controlled via a config yaml file that can be modified at runtime for real-time feedback. Custom images can be used to theme the scenario.

#### Sample Demo ####
The following demo illustrates  space-themed scenario using various actors with different behaviors:
* Blue rocket - wander around, avoiding asteroids and red ships
* White rockets - seek blue rocket, avoid asteroids, planets, and red ships
* Red ships - seek blue rocket, avoid asteroids and planets
* Asteroids - rotate, move in constant direction

![01](/demos/space-theme.gif "01")