#include "home.h"

#include "imgui_impl_dx11.h"
#include "imgui_impl_win32.h"
#include "../../globals/globals.h"
#include "../../helper/helper.h"

void homeInit(const HWND hwnd)
{
    IMGUI_CHECKVERSION();
    ImGui::CreateContext();
    ImGuiIO& io = ImGui::GetIO();
    (void)io;
    io.ConfigFlags |= ImGuiConfigFlags_NavEnableKeyboard; // Enable Keyboard Controls

    // Setup Dear ImGui style
    //ImGui::StyleColorsDark();
    ImGui::StyleColorsLight();

    // Setup Platform/Renderer backends
    ImGui_ImplWin32_Init(hwnd);
    ImGui_ImplDX11_Init(g_pd3dDevice, g_pd3dDeviceContext);
    io.Fonts->AddFontFromFileTTF("c:\\Windows\\Fonts\\segoeui.ttf", 18.0f);
}

void homeDestroy()
{
    ImGui_ImplDX11_Shutdown();
    ImGui_ImplWin32_Shutdown();
    ImGui::DestroyContext();
}

void homeUI()
{
    ImGuiWindowFlags win_flags = 0;
    win_flags |= ImGuiWindowFlags_NoTitleBar;
    win_flags |= ImGuiWindowFlags_NoMove;
    win_flags |= ImGuiWindowFlags_NoResize;
    ImGui::SetNextWindowPos(ImVec2(0, 0), ImGuiCond_FirstUseEver);
    ImGui::SetNextWindowSize(ImVec2(380, 440), ImGuiCond_Always);

    ID3D11ShaderResourceView* menu = nullptr;
    const bool menu_ret = LoadTextureFromFile("E:/cpp/gus_gui/assets/images/menu.png", &menu);
    IM_ASSERT(menu_ret);

    ID3D11ShaderResourceView* setting = nullptr;
    const bool setting_ret = LoadTextureFromFile("E:/cpp/gus_gui/assets/images/setting.png", &setting);
    IM_ASSERT(setting_ret);


    ImGui::Begin("Home view", &show_home, win_flags);
    const ImVec2 windowSize = ImGui::GetWindowSize();

    ImGui::Image(menu, ImVec2(24, 24));
    ImGui::SetItemTooltip("Menu");
    if (ImGui::IsItemHovered())
    {
        ImGui::SetMouseCursor(ImGuiMouseCursor_Hand);
    }
    if (ImGui::IsItemClicked())
    {
        show_manage = true;
    }
    const ImVec2 avail = ImGui::GetContentRegionAvail();
    ImGui::SameLine(avail.x - 24);
    ImGui::Image(setting, ImVec2(24, 24));
    ImGui::SetItemTooltip("Settings");
    if (ImGui::IsItemHovered())
    {
        ImGui::SetMouseCursor(ImGuiMouseCursor_Hand);
    }
    if (ImGui::IsItemClicked())
    {
        show_settings = true;
    }

    ImGui::Text("This is home view");

    ImGui::End();
}
