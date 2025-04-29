#define NOB_IMPLEMENTATION
#include "nob.h"

#ifdef _WIN32
    #include <windows.h>
#else
    #include <libgen.h> // For dirname()
#endif

#define BUILD_DIR "./build"
#define SRC_DIR "./src"
#define BIN_DIR "./bin"
#define RENDERER_DIR "./renderer"
#define RENDERER_LIB_A "/librenderer.a"
#define RENDERER_LIB_H "/librenderer.h"
#define OUTPUT_EXE BIN_DIR "/goui.exe"
#define INCLUDE_DIR "./renderer/include"
#define LIB_DIR "./renderer/lib/Release"

void get_root_dir(char *exe_path, char *root_dir, size_t size) {
    strncpy(root_dir, exe_path, size - 1);
    root_dir[size - 1] = '\0'; // Ensure null termination

#ifdef _WIN32
    // Windows: Use backslash (`\`) as path separator
    for (char *p = root_dir + strlen(root_dir); p >= root_dir; --p) {
        if (*p == '\\') {
            *p = '\0'; // Trim after last backslash
            break;
        }
    }
#else
    // Linux/macOS: Use `dirname()` from libgen.h
    strcpy(root_dir, dirname(root_dir));
#endif
}

int main(int argc, char **argv) {
    NOB_GO_REBUILD_URSELF(argc, argv);
    Nob_Cmd cmd = {0};

    // Check for "build-unilog" flag
    bool build_renderer = false;
    bool run = false;
    for (int i = 1; i < argc; i++) {
        if (strcmp(argv[i], "--build-renderer") == 0 || strcmp(argv[i], "-br") == 0) {
            build_renderer = true;
        }
        if (strcmp(argv[i], "--run") == 0 || strcmp(argv[i], "-r") == 0) {
            run = true;
        }
    }

    // Get root directory from argv[0]
    char root_dir[1024];
    get_root_dir(argv[0], root_dir, sizeof(root_dir));

    nob_log(NOB_INFO, "Root directory: %s\n", root_dir);

    // delete bin and build directories if they exist
    #ifdef _WIN32
    //
        nob_cmd_append(&cmd, "powershell.exe", "-Command", "Remove-Item", "-Recurse", "-Force", BUILD_DIR);
        nob_cmd_run_sync_and_reset(&cmd);
    #endif
    if (!nob_mkdir_if_not_exists(BUILD_DIR) || !nob_mkdir_if_not_exists(BIN_DIR)) return 1;
    
    if (build_renderer) {
        if (!nob_set_current_dir(RENDERER_DIR)) return 1;
        // Compile Unilog as a C archive
        //  gcc -c renderer.c -o renderer.o -Iinclude -Iexternal\glfw -Iexternal\glad -Iexternal\stb
        // gcc -c external/glad/glad.c -o glad.o -Iinclude -Iexternal\glfw -Iexternal\glad
        //  ar rcs librender.a renderer.o glad.o 
        nob_cmd_append(&cmd, "gcc", "-c", "renderer.c", "-o", "./../" BUILD_DIR "/renderer.o",  "-Iinclude", "-Iexternal/glfw", "-Iexternal/glad", "-Iexternal/stb");
        if (!nob_cmd_run_sync_and_reset(&cmd)) return 1;
        // TODO: should we save it
        nob_cmd_append(&cmd, "gcc", "-c", "external/glad/glad.c", "-o", "./../" BUILD_DIR "/glad.o", "-Iinclude", "-Iexternal/glfw", "-Iexternal/glad");
        if (!nob_cmd_run_sync_and_reset(&cmd)) return 1;
        nob_cmd_append(&cmd, "ar", "rcs", "./../" BUILD_DIR RENDERER_LIB_A, "./../" BUILD_DIR "/renderer.o", "./../" BUILD_DIR "/glad.o");
        if (!nob_cmd_run_sync_and_reset(&cmd)) return 1;
        if (!nob_set_current_dir(root_dir)) return 1;
    }
    // TODO: move inside build-renderer
    if (!nob_copy_file(BUILD_DIR RENDERER_LIB_A, LIB_DIR RENDERER_LIB_A)) return 1;

    if (!nob_set_current_dir(root_dir)) return 1;
    nob_cmd_append(&cmd, "go", "build", "-o", OUTPUT_EXE, "main.go");

    if (!nob_cmd_run_async_and_reset(&cmd)) return 1;
    nob_log(NOB_INFO, "Build completed successfully.\n");

    if(run){
        nob_cmd_append(&cmd, OUTPUT_EXE);
        if (!nob_cmd_run_async_and_reset(&cmd)) return 1;
        nob_log(NOB_INFO, "Running %s...\n", OUTPUT_EXE);
    }

    return 0;
}