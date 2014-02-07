#version 330

layout(location=0) in vec4 position;
layout(location=1) in vec4 color;

smooth out vec4 theColor;

uniform mat4 perpectiveMatrix;
uniform float loopDuration;
uniform float time;

void main() {
	theColor=color;
	float timeScale = 3.14159f * 2.0f / loopDuration;

	float currTime = mod(time, loopDuration);
	vec4 totalOffset = vec4(
			cos(currTime * timeScale) * 0.9f,
			sin(currTime * timeScale) * 0.9f,
			0.0f,
			0.0f);

	gl_Position = (position+ totalOffset) * perpectiveMatrix;
}

