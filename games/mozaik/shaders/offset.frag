#version 330

out vec4 outputColor;

uniform float fragLoopDuration;
uniform float time;

const vec4 firstColor = vec4(1.0f, 1.0f, 1.0f, 1.0f);
const vec4 secondColor = vec4(1.0f, 0.0f, 0.0f, 1.0f);

void main()
{
	float timeScale = 3.14159f *2.0f/ fragLoopDuration;

	float currTime = mod(time, fragLoopDuration);

	outputColor = mix(firstColor, secondColor, 1-abs(sin(currTime * timeScale)));
}
