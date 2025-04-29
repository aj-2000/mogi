#include <stdlib.h>
#include <stdio.h>
#include <math.h>
#include <string.h>
#include <stdarg.h> // Needed for va_list in dprintf

#include "renderer.h"

#include "external/glad/glad.h"  // Include GLAD for OpenGL function loading
#include "external/glfw/glfw3.h" // Include GLFW for window management

#define STB_TRUETYPE_IMPLEMENTATION
#include "external/stb/stb_truetype.h"

// --- Debug Printf ---
int dprintf(const char *format, ...) {
    #ifdef DEBUG
        va_list args;
        int result;
        // Use vfprintf to stderr for debug messages
        va_start(args, format);
        result = vfprintf(stderr, format, args);
        va_end(args);
        return result;
    #else
        (void)format; // Prevents "unused parameter" warning
        return 0;
    #endif
}

// --- GLFW Framebuffer Resize Callback ---
// This function gets called by GLFW whenever the window size changes
void framebuffer_size_callback(GLFWwindow* window, int width, int height) {
    if (width <= 0 || height <= 0) {
        // Avoid division by zero or invalid viewport if minimized
        return;
    }

    // 1. Update OpenGL Viewport
    glViewport(0, 0, width, height);

    // 2. Update Projection Matrix
    glMatrixMode(GL_PROJECTION);
    glLoadIdentity();
    // Use the new width and height for the orthogonal projection
    glOrtho(0.0, (double)width, (double)height, 0.0, -1.0, 1.0);
    glMatrixMode(GL_MODELVIEW);
    glLoadIdentity(); // Optional: Reset modelview matrix as well

    // 3. Update stored size in Renderer context (if needed)
    Renderer* ctx = (Renderer*)glfwGetWindowUserPointer(window);
    if (ctx) {
        ctx->current_width = width;
        ctx->current_height = height;
    }
    dprintf("Window resized to %d x %d. Viewport and Ortho updated.\n", width, height);
}

// --- Renderer Creation ---
void* create_renderer(int width, int height, const char* title) {
    if (!glfwInit()) {
        fprintf(stderr, "ERROR: Failed to initialize GLFW\n");
        return NULL;
    }

    // Optional: Add window hints before creation if needed
    // glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
    // glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 3);
    // glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE); // Requires shader-based rendering

    GLFWwindow* window = glfwCreateWindow(width, height, title, NULL, NULL);
    if (!window) {
        fprintf(stderr, "ERROR: Failed to create GLFW window\n");
        glfwTerminate();
        return NULL;
    }

    glfwMakeContextCurrent(window);

    if (!gladLoadGLLoader((GLADloadproc)glfwGetProcAddress)) {
         fprintf(stderr, "ERROR: Failed to initialize GLAD\n");
         glfwDestroyWindow(window);
         glfwTerminate();
         return NULL;
    }

    printf("OpenGL Version: %s\n", glGetString(GL_VERSION));
    printf("GLSL Version: %s\n", glGetString(GL_SHADING_LANGUAGE_VERSION));
    printf("Renderer: %s\n", glGetString(GL_RENDERER));
    printf("Vendor: %s\n", glGetString(GL_VENDOR));

    // Allocate Renderer struct
    Renderer* renderer = (Renderer*)malloc(sizeof(Renderer));
    if (!renderer) {
        fprintf(stderr, "ERROR: Failed to allocate memory for Renderer\n");
        glfwDestroyWindow(window);
        glfwTerminate();
        return NULL;
    }
    renderer->window = window;
    renderer->current_width = width;  // Store initial size
    renderer->current_height = height;

    // Store pointer to Renderer struct in GLFW window for access in callbacks
    glfwSetWindowUserPointer(window, renderer);

    // *** REGISTER THE RESIZE CALLBACK HERE ***
    glfwSetFramebufferSizeCallback(window, framebuffer_size_callback);

    // *** IMPORTANT: Call the callback ONCE manually to set initial state ***
    // This ensures viewport/projection are set correctly even if no resize happens
    // It uses the *framebuffer* size, which might differ from window size on some systems (e.g., Retina displays)
    int fb_width, fb_height;
    glfwGetFramebufferSize(window, &fb_width, &fb_height);
    framebuffer_size_callback(window, fb_width, fb_height);
    // Alternatively, you could just call the glViewport/glOrtho here directly
    // using the initial width/height, but using the callback ensures consistency.
    // glViewport(0, 0, width, height);
    // glMatrixMode(GL_PROJECTION);
    // glLoadIdentity();
    // glOrtho(0.0, (double)width, (double)height, 0.0, -1.0, 1.0);
    // glMatrixMode(GL_MODELVIEW);
    // glLoadIdentity();


    // Enable alpha blending
    glEnable(GL_BLEND);
    glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);

    // Disable depth testing for 2D
    glDisable(GL_DEPTH_TEST);

    // Set VSync (optional, 1 = on, 0 = off)
    glfwSwapInterval(1);

    dprintf("Renderer created successfully.\n");
    return renderer;
}

// --- Renderer Destruction ---
void destroy_renderer(void* renderer_ptr) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (ctx) {
        if (ctx->window) {
            glfwDestroyWindow(ctx->window);
        }
        free(ctx);
    }
    glfwTerminate(); // Terminate GLFW only after all windows are destroyed
    dprintf("Renderer destroyed.\n");
}

// --- Core Loop Functions ---
void clear_screen(void* renderer_ptr, ColorRGBA color) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window) return;

    // dprintf("Clearing screen with color: %f, %f, %f, %f\n", color.r, color.g, color.b, color.a);
    glClearColor(color.r, color.g, color.b, color.a);
    glClear(GL_COLOR_BUFFER_BIT); // Only need to clear color for 2D
}

void present_screen(void* renderer_ptr) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window) return;
    glfwSwapBuffers(ctx->window);
}

int window_should_close(void* renderer_ptr) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window) return 1; // Treat invalid renderer as should close
    return glfwWindowShouldClose(ctx->window);
}

void handle_events(void* renderer_ptr) {
    // No need for renderer context here, glfwPollEvents is global
    (void)renderer_ptr; // Mark as unused
    glfwPollEvents();
}

float get_delta_time(void* renderer_ptr) {
    // No need for renderer context here, glfwGetTime is global
    (void)renderer_ptr; // Mark as unused
    static double last_time = 0.0;
    double current_time = glfwGetTime();
    // Prevent huge delta time on first frame
    if (last_time == 0.0) {
        last_time = current_time;
        return 0.0f;
    }
    float delta_time = (float)(current_time - last_time);
    last_time = current_time;
    return delta_time;
}

void set_vsync(void* renderer_ptr, int vsync) {
    // No need for renderer context here, glfwSwapInterval affects current context
     (void)renderer_ptr; // Mark as unused
    glfwSwapInterval(vsync ? 1 : 0); // Enable or disable V-Sync (1 or 0)
}

// --- Window Size Getters ---
// Use the cached size for potential minor performance gain if called often
float get_window_width(void* renderer_ptr) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window) return 0.0f;
    // Return framebuffer width as it's used for rendering coordinates
    return (float)ctx->current_width;
    // Or query GLFW directly if you prefer:
    // int width, height;
    // glfwGetFramebufferSize(ctx->window, &width, &height);
    // return (float)width;
}

float get_window_height(void* renderer_ptr) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window) return 0.0f;
    // Return framebuffer height
    return (float)ctx->current_height;
    // Or query GLFW directly:
    // int width, height;
    // glfwGetFramebufferSize(ctx->window, &width, &height);
    // return (float)height;
}

Vec2 get_window_size(void* renderer_ptr) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    Vec2 size = {0.0f, 0.0f};
    if (!ctx || !ctx->window) return size;
    // Return framebuffer size
    size.x = (float)ctx->current_width;
    size.y = (float)ctx->current_height;
    // Or query GLFW directly:
    // int width, height;
    // glfwGetFramebufferSize(ctx->window, &width, &height);
    // size.x = (float)width;
    // size.y = (float)height;
    return size;
}


// --- Drawing Functions (Minor changes, mostly adding checks/constants) ---

// Draw a filled rectangle (RGBA)
void draw_rectangle_filled(void* renderer_ptr, Rect rect, ColorRGBA color) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window) return;
    // dprintf("Drawing filled rectangle at (%f, %f), width: %f, height: %f\n", rect.position.x, rect.position.y, rect.width, rect.height);

    glDisable(GL_TEXTURE_2D); // Ensure texturing is off
    glColor4f(color.r, color.g, color.b, color.a);

    // Use GL_QUADS (though deprecated, matches original code)
    glBegin(GL_QUADS);
    glVertex2f(rect.position.x, rect.position.y); // Top-left
    glVertex2f(rect.position.x + rect.width, rect.position.y); // Top-right
    glVertex2f(rect.position.x + rect.width, rect.position.y + rect.height); // Bottom-right
    glVertex2f(rect.position.x, rect.position.y + rect.height); // Bottom-left
    glEnd();
}

// Draw a rectangle outline (RGBA)
void draw_rectangle_outline(void* renderer_ptr, Rect rect, ColorRGBA color) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window) return;
    // dprintf("Drawing rectangle outline at (%f, %f), width: %f, height: %f\n", rect.position.x, rect.position.y, rect.width, rect.height);

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
void draw_rectangle_filled_outline(void* renderer_ptr, Rect rect, ColorRGBA fill_color, ColorRGBA outline_color) {
    // No context check needed here as called functions do it
    draw_rectangle_filled(renderer_ptr, rect, fill_color);
    draw_rectangle_outline(renderer_ptr, rect, outline_color);
}

// Draw a circle (using GL_TRIANGLE_FAN for filled)
void draw_circle_filled(void* renderer_ptr, Circle circle, ColorRGBA color) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window || circle.radius <= 0.0f) return;
    // dprintf("Drawing filled circle at (%f, %f), radius: %f\n", circle.position.x, circle.position.y, circle.radius);

    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);

    // More segments for larger circles, fewer for smaller ones?
    int num_segments = (int)(fmax(10.0f, fmin(circle.radius * 1.5f, 60.0f))); // Dynamic segments
    float angle_step = 2.0f * (float)M_PI / num_segments;

    glBegin(GL_TRIANGLE_FAN);
    glVertex2f(circle.position.x, circle.position.y); // Center
    for (int i = 0; i <= num_segments; i++) {
        float angle = i * angle_step;
        float x = circle.position.x + cosf(angle) * circle.radius;
        float y = circle.position.y + sinf(angle) * circle.radius;
        glVertex2f(x, y);
    }
    glEnd();
}

// Draw a circle outline (using GL_LINE_LOOP)
void draw_circle_outline(void* renderer_ptr, Circle circle, ColorRGBA color) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window || circle.radius <= 0.0f) return;
    // dprintf("Drawing circle outline at (%f, %f), radius: %f\n", circle.position.x, circle.position.y, circle.radius);

    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);

    int num_segments = (int)(fmax(10.0f, fmin(circle.radius * 1.5f, 60.0f)));
    float angle_step = 2.0f * (float)M_PI / num_segments;

    glBegin(GL_LINE_LOOP);
    for (int i = 0; i < num_segments; i++) { // Use < num_segments for LINE_LOOP
        float angle = i * angle_step;
        float x = circle.position.x + cosf(angle) * circle.radius;
        float y = circle.position.y + sinf(angle) * circle.radius;
        glVertex2f(x, y);
    }
    glEnd();
}

// Draw a filled circle with outline (RGBA)
void draw_circle_filled_outline(void* renderer_ptr, Circle circle, ColorRGBA fill_color, ColorRGBA outline_color) {
    draw_circle_filled(renderer_ptr, circle, fill_color);
    draw_circle_outline(renderer_ptr, circle, outline_color);
}


// Draw a thick line (using GL_QUADS)
void draw_line_thick(void* renderer_ptr, Line line, ColorRGBA color, float thickness) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window || thickness <= 0.0f) return;
    // dprintf("Drawing thick line from (%f, %f) to (%f, %f) with thickness %f\n", line.start.x, line.start.y, line.end.x, line.end.y, thickness);

    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);

    Vec2 dir = {line.end.x - line.start.x, line.end.y - line.start.y};
    float length = sqrtf(dir.x * dir.x + dir.y * dir.y);
    if (length < 0.0001f) return; // Avoid division by zero for zero-length lines

    // Normalize direction
    dir.x /= length;
    dir.y /= length;

    // Calculate perpendicular vector (half thickness)
    float half_thickness = thickness / 2.0f;
    Vec2 perp = {-dir.y * half_thickness, dir.x * half_thickness};

    glBegin(GL_QUADS);
    glVertex2f(line.start.x + perp.x, line.start.y + perp.y);
    glVertex2f(line.start.x - perp.x, line.start.y - perp.y);
    glVertex2f(line.end.x - perp.x, line.end.y - perp.y);
    glVertex2f(line.end.x + perp.x, line.end.y + perp.y);
    glEnd();
}

// Draw a dashed line (using GL_LINES)
void draw_line_dashed(void* renderer_ptr, Line line, ColorRGBA color, float dash_length, float gap_length) {
    Renderer* ctx = (Renderer*)renderer_ptr;
     if (!ctx || !ctx->window || dash_length <= 0.0f || gap_length < 0.0f) return;
    // dprintf("Drawing dashed line from (%f, %f) to (%f, %f)\n", line.start.x, line.start.y, line.end.x, line.end.y);

    glDisable(GL_TEXTURE_2D);
    glColor4f(color.r, color.g, color.b, color.a);

    Vec2 dir = {line.end.x - line.start.x, line.end.y - line.start.y};
    float total_length = sqrtf(dir.x * dir.x + dir.y * dir.y);
    if (total_length < 0.0001f) return;

    // Normalize direction
    dir.x /= total_length;
    dir.y /= total_length;

    float segment_length = dash_length + gap_length;
    if (segment_length < 0.0001f) return; // Avoid issues if dash+gap is zero

    float current_dist = 0.0f;
    glBegin(GL_LINES);
    while (current_dist < total_length) {
        float dash_end_dist = fminf(current_dist + dash_length, total_length);

        Vec2 start_pt = {line.start.x + dir.x * current_dist, line.start.y + dir.y * current_dist};
        Vec2 end_pt = {line.start.x + dir.x * dash_end_dist, line.start.y + dir.y * dash_end_dist};

        glVertex2f(start_pt.x, start_pt.y);
        glVertex2f(end_pt.x, end_pt.y);

        current_dist += segment_length;
    }
    glEnd();
}

// Draw a dotted line (using draw_circle_filled)
void draw_line_dotted(void* renderer_ptr, Line line, ColorRGBA color, float dot_radius, float gap_factor) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window || dot_radius <= 0.0f || gap_factor <= 0.0f) return;
    // dprintf("Drawing dotted line from (%f, %f) to (%f, %f)\n", line.start.x, line.start.y, line.end.x, line.end.y);

    // Note: Disabling/enabling texture inside the loop could be slow if many dots
    // Consider managing texture state outside if performance is critical.
    // glDisable(GL_TEXTURE_2D); // Done within draw_circle_filled

    Vec2 dir = {line.end.x - line.start.x, line.end.y - line.start.y};
    float total_length = sqrtf(dir.x * dir.x + dir.y * dir.y);
    if (total_length < 0.0001f) return;

    // Normalize direction
    dir.x /= total_length;
    dir.y /= total_length;

    float spacing = dot_radius * 2.0f * gap_factor; // Distance between dot centers
    if (spacing < 0.0001f) spacing = dot_radius * 2.0f; // Ensure some spacing

    int num_dots = (int)(total_length / spacing) + 1;

    for (int i = 0; i < num_dots; i++) {
        float current_dist = fminf(i * spacing, total_length); // Clamp last dot
        Vec2 dot_pos = {line.start.x + dir.x * current_dist, line.start.y + dir.y * current_dist};
        Circle dot = {dot_pos, dot_radius};
        // This will repeatedly call glDisable(GL_TEXTURE_2D)
        draw_circle_filled(renderer_ptr, dot, color);
    }
}


// --- Font Loading ---
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
    if (file_size <= 0) {
         fprintf(stderr, "ERROR: Invalid font file size for: %s\n", font_path);
         fclose(font_file);
         return NULL;
    }

    unsigned char* ttf_buffer = (unsigned char*)malloc(file_size);
    if (!ttf_buffer) {
        fprintf(stderr, "ERROR: Failed to allocate memory for font buffer (%ld bytes)\n", file_size);
        fclose(font_file);
        return NULL;
    }

    size_t read_size = fread(ttf_buffer, 1, file_size, font_file);
    fclose(font_file); // Close file ASAP
    if (read_size != file_size) {
        fprintf(stderr, "ERROR: Failed to read entire font file: %s (read %zu, expected %ld)\n", font_path, read_size, file_size);
        free(ttf_buffer);
        return NULL;
    }

    // Allocate memory for the font data structure
    FontData* font_data = (FontData*)calloc(1, sizeof(FontData)); // Use calloc to zero-initialize
    if (!font_data) {
        fprintf(stderr, "ERROR: Failed to allocate memory for FontData\n");
        free(ttf_buffer);
        return NULL;
    }
    font_data->ttf_buffer = ttf_buffer; // Store buffer pointer
    font_data->font_height_pixels = font_height_pixels;

    // Prepare temporary bitmap for stb_truetype to render into
    // Ensure atlas dimensions are reasonable
    if (FONT_ATLAS_WIDTH <= 0 || FONT_ATLAS_HEIGHT <= 0) {
         fprintf(stderr, "ERROR: Invalid FONT_ATLAS dimensions (%d x %d)\n", FONT_ATLAS_WIDTH, FONT_ATLAS_HEIGHT);
         free(font_data->ttf_buffer);
         free(font_data);
         return NULL;
    }
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

    stbtt_PackSetOversampling(&pack_context, 1, 1); // No oversampling

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

    // --- Create OpenGL Texture ---
    glGenTextures(1, &font_data->texture_id);
    glBindTexture(GL_TEXTURE_2D, font_data->texture_id);

    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR); // Linear filtering looks better for fonts
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);

    glPixelStorei(GL_UNPACK_ALIGNMENT, 1); // Crucial for single-channel textures
    glTexImage2D(
        GL_TEXTURE_2D, 0, GL_ALPHA, // Store as alpha internally
        FONT_ATLAS_WIDTH, FONT_ATLAS_HEIGHT, 0,
        GL_ALPHA, // Source format is alpha
        GL_UNSIGNED_BYTE, temp_bitmap
    );
    glPixelStorei(GL_UNPACK_ALIGNMENT, 4); // Reset to default

    glBindTexture(GL_TEXTURE_2D, 0); // Unbind

    free(temp_bitmap); // Free CPU-side bitmap

    // --- Cache Font Metrics ---
    stbtt_fontinfo info;
    if (stbtt_InitFont(&info, font_data->ttf_buffer, stbtt_GetFontOffsetForIndex(font_data->ttf_buffer, 0))) {
        int ascent_px, descent_px, lineGap_px;
        stbtt_GetFontVMetrics(&info, &ascent_px, &descent_px, &lineGap_px);
        // Calculate scale factor based on desired pixel height
        float scale = stbtt_ScaleForPixelHeight(&info, font_height_pixels);
        font_data->ascent = (float)ascent_px * scale;
        font_data->descent = (float)descent_px * scale; // Usually negative
        font_data->line_gap = (float)lineGap_px * scale;
    } else {
         fprintf(stderr, "WARNING: Failed to init font info for metrics caching: %s\n", font_path);
         // Set some defaults?
         font_data->ascent = font_height_pixels * 0.75f;
         font_data->descent = -font_height_pixels * 0.25f;
         font_data->line_gap = 0.0f;
    }


    dprintf("Font loaded: %s (Texture ID: %u, Ascent: %.2f)\n", font_path, font_data->texture_id, font_data->ascent);
    return font_data;
}

// --- Font Destruction ---
void destroy_font(FontData* font_data) {
    if (!font_data) return;

    dprintf("Destroying font (Texture ID: %u)\n", font_data->texture_id);
    // Ensure texture ID is valid before deleting
    if (font_data->texture_id > 0) {
        glDeleteTextures(1, &font_data->texture_id);
    }
    free(font_data->ttf_buffer); // Free the font file buffer
    free(font_data);             // Free the FontData struct itself
}


// --- Text Drawing ---
void draw_text(void* renderer_ptr, FontData* font_data, const char* text, Vec2 pos, ColorRGBA color) {
    Renderer* ctx = (Renderer*)renderer_ptr;
    if (!ctx || !ctx->window || !font_data || !text) return;

    // dprintf("Drawing text: '%s' at (%f, %f)\n", text, pos.x, pos.y);

    glEnable(GL_TEXTURE_2D);
    glBindTexture(GL_TEXTURE_2D, font_data->texture_id);
    glColor4f(color.r, color.g, color.b, color.a);

    float current_x = pos.x;
    // Adjust y position to account for the baseline using cached ascent
    float current_y = pos.y + font_data->ascent;

    glBegin(GL_QUADS);
    for (const char* p = text; *p; ++p) {
        // Check if character is in the packed range
        if (*p >= FONT_FIRST_CHAR && *p < FONT_FIRST_CHAR + FONT_NUM_CHARS) {
            stbtt_aligned_quad quad;
            stbtt_GetPackedQuad(
                font_data->char_data,
                FONT_ATLAS_WIDTH, FONT_ATLAS_HEIGHT,
                *p - FONT_FIRST_CHAR, // Character index
                &current_x,           // Updated by function
                &current_y,           // Updated by function (important for vertical alignment)
                &quad,
                1                     // Align to pixel grid
            );

            // Render the quad for the character
            glTexCoord2f(quad.s0, quad.t0); glVertex2f(quad.x0, quad.y0);
            glTexCoord2f(quad.s0, quad.t1); glVertex2f(quad.x0, quad.y1);
            glTexCoord2f(quad.s1, quad.t1); glVertex2f(quad.x1, quad.y1);
            glTexCoord2f(quad.s1, quad.t0); glVertex2f(quad.x1, quad.y0);
        } else {
            // Handle characters outside the range (e.g., space, tabs, unknown)
            // Get advance width for space if possible, otherwise estimate
            // Note: Getting advance requires stbtt_fontinfo, could be slow if done per char
            // For simplicity, use an estimate based on font height
             if (*p == ' ') {
                 current_x += font_data->font_height_pixels * 0.3f; // Approx space width
             } else if (*p == '\t') {
                 current_x += font_data->font_height_pixels * 0.3f * 4; // Approx tab width
             } else {
                 current_x += font_data->font_height_pixels * 0.5f; // Approx unknown char width
             }
        }
    }
    glEnd();

    glBindTexture(GL_TEXTURE_2D, 0);
    glDisable(GL_TEXTURE_2D);
}

// --- Text Width Calculation ---
float calculate_text_width(FontData* font_data, const char* text) {
    if (!font_data || !text) return 0.0f;

    float total_width = 0.0f;
    float current_x = 0.0f;
    // Use a dummy y value, it's not needed for width calculation with GetPackedQuad
    float dummy_y = 0.0f;

    for (const char* p = text; *p; ++p) {
        if (*p >= FONT_FIRST_CHAR && *p < FONT_FIRST_CHAR + FONT_NUM_CHARS) {
            stbtt_aligned_quad quad;
            // We only care about how much current_x advances
            stbtt_GetPackedQuad(
                font_data->char_data,
                FONT_ATLAS_WIDTH, FONT_ATLAS_HEIGHT,
                *p - FONT_FIRST_CHAR,
                &current_x, // This gets updated
                &dummy_y,   // This also gets updated but we ignore it
                &quad,
                0 // No pixel alignment needed for width calculation
            );
            // After GetPackedQuad, current_x is at the start of the *next* character.
            // So, the final value of current_x represents the total width.
        } else {
             // Estimate width for non-renderable characters
             if (*p == ' ') {
                 current_x += font_data->font_height_pixels * 0.3f;
             } else if (*p == '\t') {
                 current_x += font_data->font_height_pixels * 0.3f * 4;
             } else {
                 current_x += font_data->font_height_pixels * 0.5f;
             }
        }
    }
    // The final value of current_x is the total advance width
    total_width = current_x;

    // Alternative using stbtt_GetCodepointHMetrics (potentially more accurate but needs fontinfo)
    /*
    stbtt_fontinfo info;
    if (stbtt_InitFont(&info, font_data->ttf_buffer, stbtt_GetFontOffsetForIndex(font_data->ttf_buffer, 0))) {
        float scale = stbtt_ScaleForPixelHeight(&info, font_data->font_height_pixels);
        total_width = 0;
        int prev_char = 0;
        for (const char* p = text; *p; ++p) {
            int advance, left_bearing;
            stbtt_GetCodepointHMetrics(&info, *p, &advance, &left_bearing);
            total_width += advance * scale;
            if (prev_char != 0) {
                total_width += stbtt_GetCodepointKernAdvance(&info, prev_char, *p) * scale;
            }
            prev_char = *p;
        }
    } else {
        // Fallback to estimation if fontinfo fails
        total_width = strlen(text) * font_data->font_height_pixels * 0.5f; // Very rough estimate
    }
    */

    return total_width;
}
