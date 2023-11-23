#include "home.h"

#include "imgui_impl_dx11.h"
#include "imgui_impl_win32.h"
#include "../../globals/globals.h"
#include "../../helper/helper.h"
#include "../../resource/resource.h"

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

    ImGui::Begin("Home view", &show_home, win_flags);
    const ImVec2 windowSize = ImGui::GetWindowSize();
    {
        ImGui::Image(img_menu, ImVec2(24, 24));
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
        ImGui::Image(img_setting, ImVec2(24, 24));
        ImGui::SetItemTooltip("Settings");
        if (ImGui::IsItemHovered())
        {
            ImGui::SetMouseCursor(ImGuiMouseCursor_Hand);
        }
        if (ImGui::IsItemClicked())
        {
            show_settings = true;
        }
    }

    {
        const char* items[] = { "AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "IIII", "JJJJ", "KKKK", "LLLLLLL", "MMMM", "OOOOOOO" };
        static int item_current_idx = 0;
        ImGui::Text("Full-width:");
        if (ImGui::BeginListBox("##listbox 2", ImVec2(-FLT_MIN, 5 * ImGui::GetTextLineHeightWithSpacing())))
        {
            for (int n = 0; n < IM_ARRAYSIZE(items); n++)
            {
                const bool is_selected = (item_current_idx == n);
                if (ImGui::Selectable(items[n], is_selected))
                {
                    item_current_idx = n;
                }

                if (is_selected)
                {
                    ImGui::SetItemDefaultFocus();
                }
            }
            ImGui::EndListBox();
        }
    }

    // 获取头像
    /*{
        char email[] = "MyEmailAddress@example.com ";
        auto lala = GetEmailSha(email);
        ImGui::Text(lala);
        delete[] lala;
    }*/

    ImGui::Text("This is home view");

    ImGui::End();
}
