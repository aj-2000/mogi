#ifndef RENDERER_H
#define RENDERER_H

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


// Function to create and initialize the renderer (window)
void* create_renderer(int width, int height, const char* title);

// Function to clear the screen with a color (RGBA)
void clear_screen(void* renderer, ColorRGBA color);

// Function to draw a rectangle (RGBA)
void draw_rectangle(void* renderer, Rect rect, ColorRGBA color);

// Function to present the screen (swap buffers)
void present_screen(void* renderer);

// Function to check if the window should close
int window_should_close(void* renderer);

// Function to destroy the renderer (clean up)
void destroy_renderer(void* renderer);

// Function to handle events (keyboard, mouse, etc.)
void handle_events(void* renderer);

#ifdef __cplusplus
}
#endif

#endif
