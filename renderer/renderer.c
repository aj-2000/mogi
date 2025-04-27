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


// Debug dprintf function - only prints if DEBUG is defined
int dprintf(const char *format, ...) {
    #ifdef DEBUG
        va_list args;
        int result;
        
        va_start(args, format);
        result = vdprintf(format, args);
        va_end(args);
        
        return result;
    #else
        // When DEBUG is not defined, do nothing but return 0
        (void)format; // Prevents "unused parameter" warning
        return 0;
    #endif
}

// Create and initialize the renderer (window)
void* create_renderer(int width, int height, const char* title) {
    if (!glfwInit()) {
        fprintf(stderr, "ERROR: Failed to initialize GLFW\n");
        return NULL; // Initialization failed
    }

    // Request core profile if needed, but legacy works for this example
    // glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
    // glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 3);
    // glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);

    GLFWwindow* window = glfwCreateWindow(width, height, title, NULL, NULL);
    if (!window) {
        fprintf(stderr, "ERROR: Failed to create GLFW window\n");
        glfwTerminate();
        return NULL; // Window creation failed
    }

    glfwMakeContextCurrent(window);
    if (!gladLoadGLLoader((GLADloadproc)glfwGetProcAddress)) {
         fprintf(stderr, "ERROR: Failed to initialize GLAD\n");
         glfwDestroyWindow(window);
         glfwTerminate();
         return NULL;
    }

    dprintf("OpenGL Version: %s\n", glGetString(GL_VERSION));
    dprintf("GLSL Version: %s\n", glGetString(GL_SHADING_LANGUAGE_VERSION));
    dprintf("Renderer: %s\n", glGetString(GL_RENDERER));
    dprintf("Vendor: %s\n", glGetString(GL_VENDOR));


    // Setup OpenGL for 2D rendering (orthogonal projection)
    glViewport(0, 0, width, height); // Set viewport
    glMatrixMode(GL_PROJECTION);
    glLoadIdentity();
    // Origin at top-left: (0,0) top-left, (width, height) bottom-right
    glOrtho(0.0, (double)width, (double)height, 0.0, -1.0, 1.0);
    glMatrixMode(GL_MODELVIEW);
    glLoadIdentity();

    // Enable alpha blending for text transparency
    glEnable(GL_BLEND);
    glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);


    Renderer* renderer = (Renderer*)malloc(sizeof(Renderer));
    if (!renderer) {
        fprintf(stderr, "ERROR: Failed to allocate memory for Renderer\n");
        glfwDestroyWindow(window);
        glfwTerminate();
        return NULL;
    }
    renderer->window = window;
    return renderer;
}

// Clear the screen with a color (RGBA)
void clear_screen(void* renderer, ColorRGBA color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    // dprintf("Clearing screen with color: %f, %f, %f, %f\n", color.r, color.g, color.b, color.a);  // Debug log

    glClearColor(color.r, color.g, color.b, color.a);
    glClear(GL_COLOR_BUFFER_BIT);
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
        if (ctx->window) {
            glfwDestroyWindow(ctx->window);
        }
        free(ctx);
    }
    glfwTerminate();
    dprintf("Renderer destroyed.\n");
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
    dprintf("Drawing circle at (%f, %f), radius: %f\n", circle.position.x, circle.position.y, circle.radius);  // Debug log

    // Disable texturing if it was enabled for text
    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);
    int num_segments = 30; // Reduced segments for performance
    float angle_step = 2.0f * 3.1415926535f / num_segments;

    glBegin(GL_TRIANGLE_FAN);
    glVertex2f(circle.position.x, circle.position.y); // Center of the circle
    for (int i = 0; i <= num_segments; i++) {
        float angle = i * angle_step;
        float x = circle.position.x + cosf(angle) * circle.radius;
        float y = circle.position.y + sinf(angle) * circle.radius;
        glVertex2f(x, y);
    }
    glEnd();
}

// Draw a filled circle (RGBA)
void draw_circle_filled(void* renderer, Circle circle, ColorRGBA color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    dprintf("Drawing filled circle at (%f, %f), radius: %f\n", circle.position.x, circle.position.y, circle.radius);  // Debug log

    // Disable texturing if it was enabled for text
    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);
    int num_segments = 30; // Reduced segments for performance
    float angle_step = 2.0f * 3.1415926535f / num_segments;

    glBegin(GL_TRIANGLE_FAN);
    glVertex2f(circle.position.x, circle.position.y); // Center of the circle
    for (int i = 0; i <= num_segments; i++) {
        float angle = i * angle_step;
        float x = circle.position.x + cosf(angle) * circle.radius;
        float y = circle.position.y + sinf(angle) * circle.radius;
        glVertex2f(x, y);
    }
    glEnd();
}

// Draw a thick line (RGBA)
void draw_line_thick(void* renderer, Line line, ColorRGBA color, float thickness) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    // dprintf("Drawing thick line from (%f, %f) to (%f, %f) with thickness %f\n", line.start.x, line.start.y, line.end.x, line.end.y, thickness);  // Debug log

    // Disable texturing if it was enabled for text
    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);

    // Calculate the direction vector of the line
    Vec2 dir = {line.end.x - line.start.x, line.end.y - line.start.y};
    float length = sqrtf(dir.x * dir.x + dir.y * dir.y);
    if (length == 0.0f) return; // Avoid division by zero

    // Normalize the direction vector
    dir.x /= length;
    dir.y /= length;

    // Calculate perpendicular vector for thickness
    Vec2 perp = {-dir.y * thickness / 2.0f, dir.x * thickness / 2.0f};

    glBegin(GL_QUADS);
    glVertex2f(line.start.x + perp.x, line.start.y + perp.y); // Start point offset
    glVertex2f(line.start.x - perp.x, line.start.y - perp.y); // Start point offset
    glVertex2f(line.end.x - perp.x, line.end.y - perp.y);     // End point offset
    glVertex2f(line.end.x + perp.x, line.end.y + perp.y);     // End point offset
    glEnd();
}

// Draw a dashed line (RGBA)
void draw_line_dashed(void* renderer, Line line, ColorRGBA color, float dash_length, float gap_length) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    // dprintf("Drawing dashed line from (%f, %f) to (%f, %f)\n", line.start.x, line.start.y, line.end.x, line.end.y);  // Debug log

    // Disable texturing if it was enabled for text
    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);

    // Calculate the direction vector of the line
    Vec2 dir = {line.end.x - line.start.x, line.end.y - line.start.y};
    float length = sqrtf(dir.x * dir.x + dir.y * dir.y);
    if (length == 0.0f) return; // Avoid division by zero

    // Normalize the direction vector
    dir.x /= length;
    dir.y /= length;

    // Calculate the number of dashes and gaps needed
    float total_length = length;
    int num_segments = (int)(total_length / (dash_length + gap_length));
    
    for (int i = 0; i < num_segments; i++) {
        Vec2 start_offset = {line.start.x + dir.x * i * (dash_length + gap_length), line.start.y + dir.y * i * (dash_length + gap_length)};
        Vec2 end_offset = {start_offset.x + dir.x * dash_length, start_offset.y + dir.y * dash_length};

        glBegin(GL_LINES);
        glVertex2f(start_offset.x, start_offset.y);
        glVertex2f(end_offset.x, end_offset.y);
        glEnd();
    }
}

// Draw a circle with outline (RGBA)
void draw_circle_outline(void* renderer, Circle circle, ColorRGBA color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    dprintf("Drawing circle outline at (%f, %f), radius: %f\n", circle.position.x, circle.position.y, circle.radius);  // Debug log

    // Disable texturing if it was enabled for text
    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);
    int num_segments = 30; // Reduced segments for performance
    float angle_step = 2.0f * 3.1415926535f / num_segments;

    glBegin(GL_LINE_LOOP);
    for (int i = 0; i < num_segments; i++) {
        float angle = i * angle_step;
        float x = circle.position.x + cosf(angle) * circle.radius;
        float y = circle.position.y + sinf(angle) * circle.radius;
        glVertex2f(x, y);
    }
    glEnd();
}


// Draw a filled circle with outline (RGBA)
void draw_circle_filled_outline(void* renderer, Circle circle, ColorRGBA fill_color, ColorRGBA outline_color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    dprintf("Drawing filled circle with outline at (%f, %f), radius: %f\n", circle.position.x, circle.position.y, circle.radius);  // Debug log

    // Draw filled circle
    draw_circle_filled(renderer, circle, fill_color);

    // Draw outline
    draw_circle_outline(renderer, circle, outline_color);
}


// Draw a dotted line (RGBA)
void draw_line_dotted(void* renderer, Line line, ColorRGBA color, float dot_radius) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    dprintf("Drawing dotted line from (%f, %f) to (%f, %f)\n", line.start.x, line.start.y, line.end.x, line.end.y);  // Debug log

    // Disable texturing if it was enabled for text
    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);

    // Calculate the direction vector of the line
    Vec2 dir = {line.end.x - line.start.x, line.end.y - line.start.y};
    float length = sqrtf(dir.x * dir.x + dir.y * dir.y);
    if (length == 0.0f) return; // Avoid division by zero

    // Normalize the direction vector
    dir.x /= length;
    dir.y /= length;

    // Calculate the number of dots needed
    int num_dots = (int)(length / (dot_radius * 2.0f)); // Adjust for dot size

    for (int i = 0; i < num_dots; i++) {
        Vec2 offset = {line.start.x + dir.x * i * (dot_radius * 2.0f), line.start.y + dir.y * i * (dot_radius * 2.0f)};
        Circle dot = {offset, dot_radius};
        draw_circle_filled(renderer, dot, color); // Draw each dot as a filled circle
    }
}


// Draw a filled rectangle (RGBA)
void draw_rectangle_filled(void* renderer, Rect rect, ColorRGBA color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    dprintf("Drawing filled rectangle at (%f, %f), width: %f, height: %f\n", rect.position.x, rect.position.y, rect.width, rect.height);  // Debug log

    // Disable texturing if it was enabled for text
    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);

    glBegin(GL_QUADS);
    glVertex2f(rect.position.x, rect.position.y); // Top-left
    glVertex2f(rect.position.x + rect.width, rect.position.y); // Top-right
    glVertex2f(rect.position.x + rect.width, rect.position.y + rect.height); // Bottom-right
    glVertex2f(rect.position.x, rect.position.y + rect.height); // Bottom-left
    glEnd();
}

// Draw a  rectangle with outline (RGBA)
void draw_rectangle_outline(void* renderer, Rect rect, ColorRGBA color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    dprintf("Drawing rectangle outline at (%f, %f), width: %f, height: %f\n", rect.position.x, rect.position.y, rect.width, rect.height);  // Debug log

    // Disable texturing if it was enabled for text
    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);

    glBegin(GL_LINE_LOOP);
    glVertex2f(rect.position.x, rect.position.y); // Top-left
    glVertex2f(rect.position.x + rect.width, rect.position.y); // Top-right
    glVertex2f(rect.position.x + rect.width, rect.position.y + rect.height); // Bottom-right
    glVertex2f(rect.position.x, rect.position.y + rect.height); // Bottom-left
    glEnd();
}

// Draw a filled rectangle with outline (RGBA)
void draw_rectangle_filled_outline(void* renderer, Rect rect, ColorRGBA fill_color, ColorRGBA outline_color) {
    if (!renderer) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;  // Ensure the renderer is valid
    dprintf("Drawing filled rectangle with outline at (%f, %f), width: %f, height: %f\n", rect.position.x, rect.position.y, rect.width, rect.height);  // Debug log

    // Draw filled rectangle
    draw_rectangle_filled(renderer, rect, fill_color);

    // Draw outline
    draw_rectangle_outline(renderer, rect, outline_color);
}



// --- Font Loading Implementation ---
FontData* load_font(const char* font_path, float font_height_pixels) {
    // Read the font file
    FILE* font_file = fopen(font_path, "rb");
    if (!font_file) {
        fprintf(stderr, "ERROR: Failed to open font file: %s\n", font_path);
        return NULL;
    }

    fseek(font_file, 0, SEEK_END);
    long file_size = ftell(font_file);
    fseek(font_file, 0, SEEK_SET);

    unsigned char* ttf_buffer = (unsigned char*)malloc(file_size);
    if (!ttf_buffer) {
        fprintf(stderr, "ERROR: Failed to allocate memory for font buffer\n");
        fclose(font_file);
        return NULL;
    }

    size_t read_size = fread(ttf_buffer, 1, file_size, font_file);
    fclose(font_file);
    if (read_size != file_size) {
        fprintf(stderr, "ERROR: Failed to read entire font file: %s\n", font_path);
        free(ttf_buffer);
        return NULL;
    }

    // Allocate memory for the font data structure
    FontData* font_data = (FontData*)malloc(sizeof(FontData));
    if (!font_data) {
        fprintf(stderr, "ERROR: Failed to allocate memory for FontData\n");
        free(ttf_buffer);
        return NULL;
    }
    font_data->ttf_buffer = ttf_buffer; // Store buffer pointer
    font_data->font_height_pixels = font_height_pixels;

    // Prepare temporary bitmap for stb_truetype to render into
    unsigned char* temp_bitmap = (unsigned char*)malloc(FONT_ATLAS_WIDTH * FONT_ATLAS_HEIGHT);
    if (!temp_bitmap) {
        fprintf(stderr, "ERROR: Failed to allocate memory for font atlas bitmap\n");
        free(font_data->ttf_buffer);
        free(font_data);
        return NULL;
    }

    // Use stb_truetype to pack characters into the bitmap
    stbtt_pack_context pack_context;
    if (!stbtt_PackBegin(&pack_context, temp_bitmap, FONT_ATLAS_WIDTH, FONT_ATLAS_HEIGHT, 0, 1, NULL)) {
        fprintf(stderr, "ERROR: Failed to initialize stbtt_pack_context\n");
        free(temp_bitmap);
        free(font_data->ttf_buffer);
        free(font_data);
        return NULL;
    }

    stbtt_PackSetOversampling(&pack_context, 1, 1); // No oversampling for simplicity here

    // Pack the desired character range (ASCII 32-126)
    if (!stbtt_PackFontRange(&pack_context, ttf_buffer, 0, font_height_pixels, FONT_FIRST_CHAR, FONT_NUM_CHARS, font_data->char_data)) {
        fprintf(stderr, "ERROR: Failed to pack font range into atlas\n");
        stbtt_PackEnd(&pack_context);
        free(temp_bitmap);
        free(font_data->ttf_buffer);
        free(font_data);
        return NULL;
    }

    stbtt_PackEnd(&pack_context); // Finish packing

    // --- Create OpenGL Texture from the bitmap ---
    glGenTextures(1, &font_data->texture_id);
    glBindTexture(GL_TEXTURE_2D, font_data->texture_id);

    // Set texture parameters - Use GL_ALPHA since stbtt outputs grayscale alpha
    // Use GL_LINEAR for smoother scaling
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

    // Upload the pixel data - GL_ALPHA tells OpenGL it's single-channel (alpha) data
    // Modern OpenGL might prefer GL_RED, but GL_ALPHA works well with fixed-function blending
    glPixelStorei(GL_UNPACK_ALIGNMENT, 1); // Important for single-channel textures
    glTexImage2D(
        GL_TEXTURE_2D,
        0,                 // Mipmap level
        GL_ALPHA,          // Internal format (store as alpha)
        FONT_ATLAS_WIDTH,
        FONT_ATLAS_HEIGHT,
        0,                 // Border
        GL_ALPHA,          // Format of pixel data from stbtt
        GL_UNSIGNED_BYTE,  // Data type of pixel data
        temp_bitmap        // The pixel data buffer
    );
    glPixelStorei(GL_UNPACK_ALIGNMENT, 4); // Reset to default

    // Unbind texture
    glBindTexture(GL_TEXTURE_2D, 0);

    // Free the temporary bitmap buffer, we only need the OpenGL texture now
    free(temp_bitmap);

    // dprintf("Font loaded successfully: %s (Texture ID: %u)\n", font_path, font_data->texture_id);

    return font_data;
}

// --- Font Destruction Implementation ---
void destroy_font(FontData* font_data) {
    if (!font_data) return;

    dprintf("Destroying font (Texture ID: %u)\n", font_data->texture_id);
    glDeleteTextures(1, &font_data->texture_id); // Delete the OpenGL texture
    free(font_data->ttf_buffer);                 // Free the font file buffer
    free(font_data);                              // Free the FontData struct itself
}


// --- Text Drawing Implementation ---
void draw_text(void* renderer, FontData* font_data, const char* text, Vec2 pos, ColorRGBA color) {
    if (!renderer || !font_data || !text) return;

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return;

    dprintf("Drawing text: '%s' at (%f, %f)\n", text, pos.x, pos.y);  // Debug log
    
    // Enable texturing and bind the font atlas texture
    glEnable(GL_TEXTURE_2D);
    glBindTexture(GL_TEXTURE_2D, font_data->texture_id);

    // Set the text color - The texture's alpha will modulate this color's alpha
    glColor4f(color.r, color.g, color.b, color.a);

    // stbtt_aligned_quad requires current x and y to be passed by pointer
    // It calculates the quad vertices and updates x/y for the next character
    float current_x = pos.x;
    float current_y = pos.y;

    // Adjust y to baseline using font ascent
    float ascent = 0.0f;
    stbtt_fontinfo info;
    stbtt_InitFont(&info, font_data->ttf_buffer, stbtt_GetFontOffsetForIndex(font_data->ttf_buffer, 0));
    int ascent_px, descent_px, lineGap_px;
    stbtt_GetFontVMetrics(&info, &ascent_px, &descent_px, &lineGap_px);
    ascent = ascent_px * font_data->font_height_pixels / (float)(ascent_px - descent_px);
    current_y += ascent;

    glBegin(GL_QUADS); // Begin drawing quads for each character

    // Loop through each character in the string
    for (const char* p = text; *p; ++p) {
        // Check if character is in the range we packed (32-126)
        if (*p >= FONT_FIRST_CHAR && *p < FONT_FIRST_CHAR + FONT_NUM_CHARS) {
            stbtt_aligned_quad quad;

            // Get the quad geometry and texture coordinates for the character
            // Note: The y coordinate is adjusted based on the font baseline.
            // stbtt_GetPackedQuad calculates screen coords. The y-axis points down.
            stbtt_GetPackedQuad(
                font_data->char_data, // Character data array
                FONT_ATLAS_WIDTH,     // Atlas width
                FONT_ATLAS_HEIGHT,    // Atlas height
                *p - FONT_FIRST_CHAR, // Character index (offset from first char)
                &current_x,           // Pointer to current x position (updated by function)
                &current_y,           // Pointer to current y position (updated by function)
                &quad,                // Output quad structure
                1                     // Align to pixel grid (1 = true)
            );

            // Render the quad
            glTexCoord2f(quad.s0, quad.t0); glVertex2f(quad.x0, quad.y0); // Top-left
            glTexCoord2f(quad.s0, quad.t1); glVertex2f(quad.x0, quad.y1); // Bottom-left
            glTexCoord2f(quad.s1, quad.t1); glVertex2f(quad.x1, quad.y1); // Bottom-right
            glTexCoord2f(quad.s1, quad.t0); glVertex2f(quad.x1, quad.y0); // Top-right

        } else {
             // Handle characters outside the range (e.g., skip, draw '?')
             // For simplicity, just advance position roughly for unknown chars like space
             if (*p == ' ') {
                 current_x += font_data->font_height_pixels * 0.3f; // Approximate space width
             } else {
                 // Maybe draw a '?' or skip
                 // For now, just advance a bit
                 current_x += font_data->font_height_pixels * 0.5f;
             }
        }
    }

    glEnd(); // End drawing quads

    // Unbind texture and disable texturing
    glBindTexture(GL_TEXTURE_2D, 0);
    glDisable(GL_TEXTURE_2D);
}

float get_delta_time(void* renderer) {
    if (!renderer) return 0.0f; 

    Renderer* ctx = (Renderer*)renderer;
    if (!ctx || !ctx->window) return 0.0f;  // Ensure the renderer is valid
    
    static double last_time = 0.0;
    double current_time = glfwGetTime();
    float delta_time = (float)(current_time - last_time);
    last_time = current_time;
    return delta_time;
}
