#include <stdlib.h>
#include <stdio.h>
#include "renderer.h"
#include "external/glad/glad.h"  // Include GLAD for OpenGL function loading
#include "external/glfw/glfw3.h" // Include GLFW for window management

typedef struct {
    GLFWwindow* window;
} Renderer;

// Create and initialize the renderer (window)
void* create_renderer(int width, int height, const char* title) {
    if (!glfwInit()) {
        return NULL; // Initialization failed
    }

    GLFWwindow* window = glfwCreateWindow(width, height, title, NULL, NULL);
    if (!window) {
        glfwTerminate();
        return NULL; // Window creation failed
    }

    glfwMakeContextCurrent(window);
    gladLoadGLLoader((GLADloadproc)glfwGetProcAddress);

    // Setup OpenGL for 2D rendering (orthogonal projection)
    glMatrixMode(GL_PROJECTION);
    glLoadIdentity();
    glOrtho(0.0, width, height, 0.0, -1.0, 1.0);  // 2D orthogonal projection matrix
    glMatrixMode(GL_MODELVIEW);
    glLoadIdentity();

    Renderer* renderer = (Renderer*)malloc(sizeof(Renderer));
    renderer->window = window;
    return renderer;
}

// Clear the screen with a color (RGBA)
void clear_screen(void* renderer, ColorRGBA color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    printf("Clearing screen with color: %f, %f, %f, %f\n", color.r, color.g, color.b, color.a);  // Debug log

    glClearColor(color.r, color.g, color.b, color.a);
    glClear(GL_COLOR_BUFFER_BIT);
}

// Draw a rectangle (RGBA)
void draw_rectangle(void* renderer, Rect rect, ColorRGBA color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    printf("Drawing rectangle at (%f, %f), width: %f, height: %f\n", rect.position.x, rect.position.y, rect.width, rect.height);  // Debug log

    glColor4f(color.r, color.g, color.b, color.a);

    glBegin(GL_QUADS);
    glVertex2f(rect.position.x, rect.position.y); // Bottom-left
    glVertex2f(rect.position.x + rect.width, rect.position.y); // Bottom-right
    glVertex2f(rect.position.x + rect.width, rect.position.y + rect.height); // Top-right
    glVertex2f(rect.position.x, rect.position.y + rect.height); // Top-left
    glEnd();
}

// Present the screen (swap buffers)
void present_screen(void* renderer) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    glfwSwapBuffers(ctx->window);
}

// Check if the window should close
int window_should_close(void* renderer) {
    if (!renderer) return 1;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return 1;  // Ensure the renderer is valid
    return glfwWindowShouldClose(ctx->window);
}

// Destroy the renderer (clean up)
void destroy_renderer(void* renderer) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (ctx) {
        glfwDestroyWindow(ctx->window);
        free(ctx);
    }
    glfwTerminate();
}

void handle_events(void* renderer) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid

    glfwPollEvents();  // Poll for and process events
}