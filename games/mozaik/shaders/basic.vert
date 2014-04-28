#version 330

layout(location=0) in vec4 position;
layout(location=1) in vec4 color;

smooth out vec4 theColor;

uniform mat4 projectionView;
uniform mat4 modelView;

void main() {
	theColor = color;
	//gl_Position = (position+ totalOffset) * perpectiveMatrix;
	gl_Position = position*modelView;
}

