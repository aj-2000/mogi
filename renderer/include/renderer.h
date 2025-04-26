#ifndef RENDERER_H
#define RENDERER_H

#ifdef __cplusplus
extern "C" {
#endif

// Function to create and initialize the renderer (window)
void* create_renderer(int width, int height, const char* title);

// Function to clear the screen with a color (RGBA)
void clear_screen(void* renderer, float r, float g, float b, float a);

// Function to draw a rectangle (RGBA)
void draw_rectangle(void* renderer, float x, float y, float width, float height, float r, float g, float b, float a);

// Function to present the screen (swap buffers)
void present_screen(void* renderer);

// Function to check if the window should close
int window_should_close(void* renderer);

// Function to destroy the renderer (clean up)
void destroy_renderer(void* renderer);

#ifdef __cplusplus
}
#endif

#endif
