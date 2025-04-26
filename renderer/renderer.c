#include <stdlib.h>
#include <stdio.h>
#include <math.h>
#include <string.h>

#include "renderer.h"

#include "external/glad/glad.h"  // Include GLAD for OpenGL function loading
#include "external/glfw/glfw3.h" // Include GLFW for window management
#define STB_TRUETYPE_IMPLEMENTATION
#include "external/stb/stb_truetype.h" 


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

// Handle events (poll for and process events)
void handle_events(void* renderer) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid

    glfwPollEvents();  // Poll for and process events
}

// Draw a circle (RGBA)
void draw_circle(void* renderer, Circle circle, ColorRGBA color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    printf("Drawing circle at (%f, %f), radius: %f\n", circle.position.x, circle.position.y, circle.radius);  // Debug log

    glColor4f(color.r, color.g, color.b, color.a);
    int num_segments = 100; // Number of segments for the circle
    float angle_step = 2.0f * 3.1415926f / num_segments;

    glBegin(GL_TRIANGLE_FAN);
    glVertex2f(circle.position.x, circle.position.y); // Center of the circle
    for (int i = 0; i <= num_segments; i++) {
        float angle = i * angle_step;
        float x = circle.position.x + cos(angle) * circle.radius;
        float y = circle.position.y + sin(angle) * circle.radius;
        glVertex2f(x, y);
    }
    glEnd();
}

// Draw a line (RGBA)
void draw_line(void* renderer, Line line, ColorRGBA color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    printf("Drawing line from (%f, %f) to (%f, %f)\n", line.start.x, line.start.y, line.end.x, line.end.y);  // Debug log

    glColor4f(color.r, color.g, color.b, color.a);
    glBegin(GL_LINES);
    glVertex2f(line.start.x, line.start.y);
    glVertex2f(line.end.x, line.end.y);
    glEnd();
}

// Initialize font
FontData* init_font(const char* font_path, float pixel_height) {
    FontData* font_data = (FontData*)malloc(sizeof(FontData));
    if (!font_data) return NULL;
    
    // Load font file
    FILE* font_file = fopen(font_path, "rb");
    if (!font_file) {
        free(font_data);
        return NULL;
    }
    
    // Get file size
    fseek(font_file, 0, SEEK_END);
    long file_size = ftell(font_file);
    fseek(font_file, 0, SEEK_SET);
    
    // Read font into buffer
    font_data->font_buffer = (unsigned char*)malloc(file_size);
    if (!font_data->font_buffer) {
        fclose(font_file);
        free(font_data);
        return NULL;
    }
    
    fread(font_data->font_buffer, 1, file_size, font_file);
    fclose(font_file);
    
    // Initialize font
    if (!stbtt_InitFont(&font_data->font_info, font_data->font_buffer, 0)) {
        free(font_data->font_buffer);
        free(font_data);
        return NULL;
    }
    
    // Calculate font scaling
    font_data->scale = stbtt_ScaleForPixelHeight(&font_data->font_info, pixel_height);
    
    // Get font metrics
    stbtt_GetFontVMetrics(&font_data->font_info, &font_data->ascent, 
                          &font_data->descent, &font_data->line_gap);
    
    // Create bitmap for text rendering
    font_data->bitmap_width = 512;  // Can be adjusted based on needs
    font_data->bitmap_height = 512;
    font_data->bitmap = (unsigned char*)malloc(font_data->bitmap_width * font_data->bitmap_height);
    
    // Create texture for rendering
    glGenTextures(1, &font_data->texture_id);
    
    return font_data;
}

// Draw text using stb_truetype
void draw_text(void* renderer, const char* text, Vec2 position, ColorRGBA color, FontData* font) {
    if (!renderer || !font) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    printf("Drawing text: '%s' at (%f, %f)\n", text, position.x, position.y);  // Debug log

    // Set text color
    glColor4f(color.r, color.g, color.b, color.a);
    
    // Clear the bitmap
    memset(font->bitmap, 0, font->bitmap_width * font->bitmap_height);
    
    // Current position for rendering
    float x = 0;
    float y = 0;
    
    // Render each character
    for (const char* c = text; *c; c++) {
        // Get character bounding box
        int x0, y0, x1, y1;
        stbtt_GetCodepointBitmapBox(&font->font_info, *c, font->scale, font->scale, 
                                    &x0, &y0, &x1, &y1);
        
        // Get character advance and kerning
        int advance, lsb;
        stbtt_GetCodepointHMetrics(&font->font_info, *c, &advance, &lsb);
        
        // Calculate position in the bitmap
        int bitmap_x = (int)x + x0;
        int bitmap_y = (int)y + y0 + (font->ascent * font->scale);
        
        // Render the character to the bitmap
        stbtt_MakeCodepointBitmap(&font->font_info, 
                                 font->bitmap + bitmap_x + bitmap_y * font->bitmap_width,
                                 x1 - x0, y1 - y0, font->bitmap_width, 
                                 font->scale, font->scale, *c);
        
        // Advance position
        x += advance * font->scale;
        
        // Add kerning with the next character
        if (*(c+1)) {
            x += stbtt_GetCodepointKernAdvance(&font->font_info, *c, *(c+1)) * font->scale;
        }
    }
    
    // Bind the texture
    glBindTexture(GL_TEXTURE_2D, font->texture_id);
    
    // Update texture with the rendered text bitmap
    glTexImage2D(GL_TEXTURE_2D, 0, GL_ALPHA, font->bitmap_width, font->bitmap_height, 
                0, GL_ALPHA, GL_UNSIGNED_BYTE, font->bitmap);
    
    // Set texture parameters
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
    
    // Enable texturing and blending for text
    glEnable(GL_TEXTURE_2D);
    glEnable(GL_BLEND);
    glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);
    
    // Draw the textured quad
    float text_width = x;
    float text_height = (font->ascent - font->descent) * font->scale;
    
    glBegin(GL_QUADS);
    glTexCoord2f(0, 0); glVertex2f(position.x, position.y);
    glTexCoord2f(1, 0); glVertex2f(position.x + text_width, position.y);
    glTexCoord2f(1, 1); glVertex2f(position.x + text_width, position.y + text_height);
    glTexCoord2f(0, 1); glVertex2f(position.x, position.y + text_height);
    glEnd();
    
    // Disable texturing and blending
    glDisable(GL_BLEND);
    glDisable(GL_TEXTURE_2D);
}

// Free font resources
void destroy_font(FontData* font) {
    if (!font) return;
    
    // Free resources
    if (font->bitmap) free(font->bitmap);
    if (font->font_buffer) free(font->font_buffer);
    
    // Delete OpenGL texture
    glDeleteTextures(1, &font->texture_id);
    
    free(font);
}