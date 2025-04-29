#ifndef RENDERER_H
#define RENDERER_H

#include "../external/stb/stb_truetype.h" 
#include "../external/glad/glad.h"  // Include GLAD for OpenGL function loading
#include "../external/glfw/glfw3.h" // Include GLFW for window management

#ifdef __cplusplus
extern "C" {
#endif
typedef struct {
    float x;
    float y;
} Vec2;

typedef struct {
    float r; // 0.0f to 1.0f
    float g; // 0.0f to 1.0f
    float b; // 0.0f to 1.0f
    float a; // 0.0f to 1.0f
} ColorRGBA;

typedef struct {
    Vec2 position;
    float width;
    float height;
} Rect;

typedef struct {
    Vec2 position;
    float radius;
} Circle;

typedef struct {
    Vec2 start;
    Vec2 end;
} Line;

// --- New Font Data Structure ---
#define FONT_ATLAS_WIDTH 512
#define FONT_ATLAS_HEIGHT 512
#define FONT_FIRST_CHAR 32 // Start rendering from space character
#define FONT_NUM_CHARS 95  // Render ASCII 32 through 126

typedef struct {
    stbtt_packedchar char_data[FONT_NUM_CHARS]; // Character metrics and texture coords
    unsigned int texture_id;                    // OpenGL texture ID for the atlas
    float font_height_pixels;                   // Font height used for packing
    // Keep font buffer alive as stbtt needs it
    unsigned char* ttf_buffer;
} FontData;


// --- Function Prototypes ---

// Renderer lifecycle
void* create_renderer(int width, int height, const char* title);
void destroy_renderer(void* renderer);
int window_should_close(void* renderer);
void handle_events(void* renderer);

// Drawing operations
void clear_screen(void* renderer, ColorRGBA color);
void present_screen(void* renderer);

void draw_rectangle_outline(void* renderer, Rect rect, ColorRGBA color);
void draw_rectangle_filled(void* renderer, Rect rect, ColorRGBA color);
void draw_rectangle_filled_outline(void* renderer, Rect rect, ColorRGBA fill_color, ColorRGBA outline_color);

void draw_circle_outline(void* renderer, Circle circle, ColorRGBA color);
void draw_circle_filled(void* renderer, Circle circle, ColorRGBA color);
void draw_circle_filled_outline(void* renderer, Circle circle, ColorRGBA fill_color, ColorRGBA outline_color);

void draw_line(void* renderer, Line line, ColorRGBA color);
void draw_line_thick(void* renderer, Line line, ColorRGBA color, float thickness);
void draw_line_dashed(void* renderer, Line line, ColorRGBA color, float dash_length, float gap_length);
void draw_line_dotted(void* renderer, Line line, ColorRGBA color, float dot_radius);

// --- New Font Functions ---
FontData* load_font(const char* font_path, float font_height_pixels);
void destroy_font(FontData* font_data);
void draw_text(void* renderer, FontData* font_data, const char* text, Vec2 pos, ColorRGBA color);


// get delta time
float get_delta_time(void* renderer);
void set_vsync(void* renderer, int vsync);

Vec2 get_window_size(void* renderer);
float calculate_text_width(FontData* font_data, const char* text);

#ifdef __cplusplus
}
#endif
#endif // RENDERER_H