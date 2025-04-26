#ifndef RENDERER_H
#define RENDERER_H

#include <stdlib.h>
#include <stdio.h>
#include <math.h>
#include "renderer.h"
#include "../external/stb/stb_truetype.h" 
#include "../external/glad/glad.h"  // Include GLAD for OpenGL function loading
#include "../external/glfw/glfw3.h" // Include GLFW for window management


#ifdef __cplusplus
extern "C" {
#endif
// Base vector types
typedef struct {
    float x, y;
} Vec2;

typedef struct {
    float x, y, z;
} Vec3;

typedef struct {
    float x, y, z, w;
} Vec4;

// Color type
typedef struct {
    float r, g, b, a;
} ColorRGBA;

// Shapes built using the base types
typedef struct {
    Vec2 position;
    float width, height;
} Rect;

typedef struct {
    Vec2 position;
    float radius;
} Circle;

typedef struct {
    Vec2 start;
    Vec2 end;
} Line;

// Font data structure for text rendering
typedef struct {
    stbtt_fontinfo font_info;
    unsigned char* font_buffer;
    float scale;
    int ascent, descent, line_gap;
    unsigned char* bitmap;
    int bitmap_width, bitmap_height;
    GLuint texture_id;
} FontData;

// Function to create and initialize the renderer (window)
void* create_renderer(int width, int height, const char* title);

// Function to clear the screen with a color (RGBA)
void clear_screen(void* renderer, ColorRGBA color);

// Function to draw a rectangle (RGBA)
void draw_rectangle(void* renderer, Rect rect, ColorRGBA color);

// Function to draw a circle (RGBA)
void draw_circle(void* renderer, Circle circle, ColorRGBA color);

// Function to draw a line (RGBA)
void draw_line(void* renderer, Line line, ColorRGBA color);

// Function to draw text (RGBA)
void draw_text(void* renderer, const char* text, Vec2 position, ColorRGBA color, FontData* font);

// Function to present the screen (swap buffers)
void present_screen(void* renderer);

// Function to check if the window should close
int window_should_close(void* renderer);

// Function to destroy the renderer (clean up)
void destroy_renderer(void* renderer);

// Function to handle events (keyboard, mouse, etc.)
void handle_events(void* renderer);

// Function to load a font and initialize font data
FontData* init_font(const char* font_path, float pixel_height);

// Free font resources
void destroy_font(FontData* font);


#ifdef __cplusplus
}
#endif

#endif
